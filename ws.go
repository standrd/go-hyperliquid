package hyperliquid

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// pingInterval is the interval for sending ping messages to keep WebSocket alive
	pingInterval = 50 * time.Second
	// gracefulCloseTimeout is the timeout for graceful WebSocket close
	gracefulCloseTimeout = 10 * time.Second
)

type WebsocketClient struct {
	url                 string
	conn                *websocket.Conn
	mu                  sync.RWMutex
	writeMu             sync.Mutex
	subscriptions       map[subKey]map[int]*subscriptionCallback
	parsedSubscriptions map[subKey]map[int]*parsedSubscriptionCallback[any]
	// Reverse lookup maps for efficient unsubscribe by ID
	subscriptionByID       map[int]*subscriptionInfo
	parsedSubscriptionByID map[int]*parsedSubscriptionInfo
	nextSubID              atomic.Int32
	done                   chan struct{}
	reconnectWait          time.Duration
}

type subscriptionInfo struct {
	key subKey
	sub Subscription
}

type parsedSubscriptionInfo struct {
	key subKey
	sub ParsedSubscriptionInterface
}

func NewWebsocketClient(baseURL string) *WebsocketClient {
	if baseURL == "" {
		baseURL = MainnetAPIURL
	}
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("invalid URL: %v", err)
	}
	parsedURL.Scheme = "wss"
	parsedURL.Path = "/ws"
	wsURL := parsedURL.String()

	return &WebsocketClient{
		url:                    wsURL,
		subscriptions:          make(map[subKey]map[int]*subscriptionCallback),
		parsedSubscriptions:    make(map[subKey]map[int]*parsedSubscriptionCallback[any]),
		subscriptionByID:       make(map[int]*subscriptionInfo),
		parsedSubscriptionByID: make(map[int]*parsedSubscriptionInfo),
		done:                   make(chan struct{}),
		reconnectWait:          time.Second,
	}
}

func (w *WebsocketClient) Connect(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		return nil
	}

	dialer := websocket.Dialer{}

	//nolint:bodyclose // WebSocket connections don't have response bodies to close
	conn, _, err := dialer.DialContext(ctx, w.url, nil)
	if err != nil {
		return fmt.Errorf("websocket dial: %w", err)
	}

	w.conn = conn

	go w.readPump(ctx)
	go w.pingPump(ctx)

	return w.resubscribeAll()
}

func (w *WebsocketClient) Subscribe(sub Subscription, callback func(WSMessage)) (int, error) {
	if callback == nil {
		return 0, fmt.Errorf("callback cannot be nil")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	key := sub.key()
	id := int(w.nextSubID.Add(1))

	if w.subscriptions[key] == nil {
		w.subscriptions[key] = make(map[int]*subscriptionCallback)
	}

	w.subscriptions[key][id] = &subscriptionCallback{
		id:       id,
		callback: callback,
	}

	// Store reverse lookup
	w.subscriptionByID[id] = &subscriptionInfo{
		key: key,
		sub: sub,
	}

	if err := w.sendSubscribe(sub); err != nil {
		delete(w.subscriptions[key], id)
		delete(w.subscriptionByID, id)
		return 0, fmt.Errorf("subscribe: %w", err)
	}

	return id, nil
}

// SubscribeParsed provides type-safe subscription with generic message types
func (w *WebsocketClient) SubscribeParsed(sub ParsedSubscriptionInterface) (int, error) {
	if sub == nil || sub.GetCallback() == nil {
		return 0, fmt.Errorf("subscription and callback cannot be nil")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	key := sub.GetSubscription().key()
	id := int(w.nextSubID.Add(1))

	if w.parsedSubscriptions[key] == nil {
		w.parsedSubscriptions[key] = make(map[int]*parsedSubscriptionCallback[any])
	}

	w.parsedSubscriptions[key][id] = &parsedSubscriptionCallback[any]{
		id:        id,
		callback:  sub.GetCallback(),
		unmarshal: sub.GetUnmarshaler(),
	}

	// Store reverse lookup
	w.parsedSubscriptionByID[id] = &parsedSubscriptionInfo{
		key: key,
		sub: sub,
	}

	if err := w.sendSubscribe(sub.GetSubscription()); err != nil {
		delete(w.parsedSubscriptions[key], id)
		delete(w.parsedSubscriptionByID, id)
		return 0, fmt.Errorf("subscribe: %w", err)
	}

	return id, nil
}

// Unsubscribe unsubscribes from a subscription by ID
func (w *WebsocketClient) Unsubscribe(id int) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if it's a regular subscription
	if info, ok := w.subscriptionByID[id]; ok {
		subs := w.subscriptions[info.key]
		if _, ok := subs[id]; !ok {
			return fmt.Errorf("subscription ID not found")
		}

		delete(subs, id)
		delete(w.subscriptionByID, id)

		if len(subs) == 0 {
			delete(w.subscriptions, info.key)
			if err := w.sendUnsubscribe(info.sub); err != nil {
				return fmt.Errorf("unsubscribe: %w", err)
			}
		}

		return nil
	}

	// Check if it's a parsed subscription
	if info, ok := w.parsedSubscriptionByID[id]; ok {
		subs := w.parsedSubscriptions[info.key]
		if _, ok := subs[id]; !ok {
			return fmt.Errorf("subscription ID not found")
		}

		delete(subs, id)
		delete(w.parsedSubscriptionByID, id)

		if len(subs) == 0 {
			delete(w.parsedSubscriptions, info.key)
			if err := w.sendUnsubscribe(info.sub.GetSubscription()); err != nil {
				return fmt.Errorf("unsubscribe: %w", err)
			}
		}

		return nil
	}

	return fmt.Errorf("subscription ID not found")
}

func (w *WebsocketClient) Close() error {
	close(w.done)

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}

// Private methods

func (w *WebsocketClient) readPump(ctx context.Context) {
	defer func() {
		w.mu.Lock()
		if w.conn != nil {
			_ = w.conn.Close() // Ignore close error in defer
			w.conn = nil
		}
		w.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.done:
			return
		default:
			_, msg, err := w.conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("websocket read error: %v", err)
					w.reconnect()
				}
				return
			}

			if string(msg) == "Websocket connection established." {
				continue
			}

			var wsMsg WSMessage
			if err := json.Unmarshal(msg, &wsMsg); err != nil {
				log.Printf("websocket message parse error: %v", err)
				continue
			}

			w.dispatch(wsMsg)
		}
	}
}

func (w *WebsocketClient) pingPump(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.done:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.sendPing(); err != nil {
				log.Printf("ping error: %v", err)
				w.reconnect()
				return
			}
		}
	}
}

func (w *WebsocketClient) dispatch(msg WSMessage) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// Handle regular subscriptions
	for key, subs := range w.subscriptions {
		// For regular subscriptions, we can match on channel type
		if key.typ == msg.Channel {
			for _, sub := range subs {
				sub.callback(msg)
			}
		}
	}

	// Handle parsed subscriptions - parse first, then match
	for key, subs := range w.parsedSubscriptions {
		// First check if the channel matches the subscription type
		if key.typ == msg.Channel {
			for _, sub := range subs {
				if parsedData, err := sub.unmarshal(msg.Data); err == nil {
					// Now that we have parsed data, check if it matches the subscription filters
					if w.matchesParsedSubscription(key, parsedData) {
						// Call the callback function directly
						if sub.callback != nil {
							callbackValue := reflect.ValueOf(sub.callback)
							if callbackValue.IsValid() && !callbackValue.IsNil() {
								callbackValue.Call([]reflect.Value{reflect.ValueOf(parsedData)})
							}
						}
					}
				} else {
					log.Printf("failed to unmarshal parsed subscription data: %v", err)
				}
			}
		}
	}
}

func (w *WebsocketClient) reconnect() {
	for {
		select {
		case <-w.done:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), gracefulCloseTimeout)
			err := w.Connect(ctx)
			cancel()
			if err == nil {
				return
			}
			time.Sleep(w.reconnectWait)
			w.reconnectWait *= 2
			if w.reconnectWait > time.Minute {
				w.reconnectWait = time.Minute
			}
		}
	}
}

func (w *WebsocketClient) resubscribeAll() error {
	for key, subs := range w.subscriptions {
		if len(subs) > 0 {
			sub := Subscription{
				Type:     key.typ,
				Coin:     key.coin,
				User:     key.user,
				Interval: key.interval,
			}
			if err := w.sendSubscribe(sub); err != nil {
				return fmt.Errorf("resubscribe: %w", err)
			}
		}
	}

	for key, subs := range w.parsedSubscriptions {
		if len(subs) > 0 {
			sub := Subscription{
				Type:     key.typ,
				Coin:     key.coin,
				User:     key.user,
				Interval: key.interval,
			}
			if err := w.sendSubscribe(sub); err != nil {
				return fmt.Errorf("resubscribe parsed: %w", err)
			}
		}
	}
	return nil
}

func (w *WebsocketClient) sendSubscribe(sub Subscription) error {
	return w.writeJSON(WsCommand{
		Method:       "subscribe",
		Subscription: &sub,
	})
}

func (w *WebsocketClient) sendUnsubscribe(sub Subscription) error {
	return w.writeJSON(WsCommand{
		Method:       "unsubscribe",
		Subscription: &sub,
	})
}

func (w *WebsocketClient) sendPing() error {
	return w.writeJSON(WsCommand{Method: "ping"})
}

func (w *WebsocketClient) writeJSON(v any) error {
	w.writeMu.Lock()
	defer w.writeMu.Unlock()

	if w.conn == nil {
		return fmt.Errorf("connection closed")
	}

	return w.conn.WriteJSON(v)
}

func (w *WebsocketClient) SubscribeToTrades(coin string, callback func(WSMessage)) (int, error) {
	sub := Subscription{Type: "trades", Coin: coin}
	return w.Subscribe(sub, callback)
}

func (w *WebsocketClient) SubscribeToOrderbook(coin string, callback func(WSMessage)) (int, error) {
	sub := Subscription{Type: "l2Book", Coin: coin}
	return w.Subscribe(sub, callback)
}

// SubscribeToAllMids subscribes to all mid prices
func (w *WebsocketClient) SubscribeToAllMids(callback func(WSMessage)) (int, error) {
	sub := Subscription{Type: "allMids"}
	return w.Subscribe(sub, callback)
}

// SubscribeToUserEvents subscribes to user events
func (w *WebsocketClient) SubscribeToUserEvents(
	user string,
	callback func(WSMessage),
) (int, error) {
	sub := Subscription{Type: "userEvents", User: user}
	return w.Subscribe(sub, callback)
}

// SubscribeToUserFills subscribes to user fills
func (w *WebsocketClient) SubscribeToUserFills(user string, callback func(WSMessage)) (int, error) {
	sub := Subscription{Type: "userFills", User: user}
	return w.Subscribe(sub, callback)
}

// SubscribeToCandles subscribes to candle data
func (w *WebsocketClient) SubscribeToCandles(
	coin, interval string,
	callback func(WSMessage),
) (int, error) {
	sub := Subscription{Type: "candle", Coin: coin, Interval: interval}
	return w.Subscribe(sub, callback)
}

// SubscribeToOrderUpdates subscribes to order updates
func (w *WebsocketClient) SubscribeToOrderUpdates(callback func(WSMessage)) (int, error) {
	sub := Subscription{Type: "orderUpdates"}
	return w.Subscribe(sub, callback)
}

// SubscribeToUserFundings subscribes to user funding updates
func (w *WebsocketClient) SubscribeToUserFundings(
	user string,
	callback func(WSMessage),
) (int, error) {
	sub := Subscription{Type: "userFundings", User: user}
	return w.Subscribe(sub, callback)
}

// SubscribeToUserNonFundingLedgerUpdates subscribes to user non-funding ledger updates
func (w *WebsocketClient) SubscribeToUserNonFundingLedgerUpdates(
	user string,
	callback func(WSMessage),
) (int, error) {
	sub := Subscription{Type: "userNonFundingLedgerUpdates", User: user}
	return w.Subscribe(sub, callback)
}

// SubscribeToWebData2 subscribes to web data v2
func (w *WebsocketClient) SubscribeToWebData2(user string, callback func(WSMessage)) (int, error) {
	sub := Subscription{Type: "webData2", User: user}
	return w.Subscribe(sub, callback)
}

// SubscribeToBBO subscribes to best bid/offer data
func (w *WebsocketClient) SubscribeToBBO(coin string, callback func(WSMessage)) (int, error) {
	sub := Subscription{Type: "bbo", Coin: coin}
	return w.Subscribe(sub, callback)
}

// SubscribeToActiveAssetCtx subscribes to active asset context
func (w *WebsocketClient) SubscribeToActiveAssetCtx(
	coin string,
	callback func(WSMessage),
) (int, error) {
	sub := Subscription{Type: "activeAssetCtx", Coin: coin}
	return w.Subscribe(sub, callback)
}

// Type-safe subscription convenience methods

// SubscribeToTradesParsed subscribes to trades with parsed Trade messages
func (w *WebsocketClient) SubscribeToTradesParsed(coin string, callback func([]WsTrade)) (int, error) {
	sub := NewParsedSubscription("trades", callback).WithCoin(coin)
	return w.SubscribeParsed(sub)
}

// SubscribeToOrderbookParsed subscribes to orderbook with parsed L2Book messages
func (w *WebsocketClient) SubscribeToOrderbookParsed(coin string, callback func(WsBook)) (int, error) {
	sub := NewParsedSubscription("l2Book", callback).WithCoin(coin)
	return w.SubscribeParsed(sub)
}

// SubscribeToAllMidsParsed subscribes to all mid prices with parsed messages
func (w *WebsocketClient) SubscribeToAllMidsParsed(callback func(map[string]string)) (int, error) {
	sub := NewParsedSubscription("allMids", callback)
	return w.SubscribeParsed(sub)
}

// SubscribeToUserEventsParsed subscribes to user events with parsed messages
func (w *WebsocketClient) SubscribeToUserEventsParsed(user string, callback func(WsUserEvent)) (int, error) {
	sub := NewParsedSubscription("userEvents", callback).WithUser(user)
	return w.SubscribeParsed(sub)
}

// SubscribeToUserFillsParsed subscribes to user fills with parsed messages
func (w *WebsocketClient) SubscribeToUserFillsParsed(user string, callback func([]WsFill)) (int, error) {
	sub := NewParsedSubscription("userFills", callback).WithUser(user)
	return w.SubscribeParsed(sub)
}

// SubscribeToCandlesParsed subscribes to candle data with parsed messages
func (w *WebsocketClient) SubscribeToCandlesParsed(coin, interval string, callback func(WsCandle)) (int, error) {
	sub := NewParsedSubscription("candle", callback).WithCoin(coin).WithInterval(interval)
	return w.SubscribeParsed(sub)
}

// SubscribeToOrderUpdatesParsed subscribes to order updates with parsed messages
func (w *WebsocketClient) SubscribeToOrderUpdatesParsed(callback func([]WsOrder)) (int, error) {
	sub := NewParsedSubscription("orderUpdates", callback)
	return w.SubscribeParsed(sub)
}

// SubscribeToActiveAssetCtxParsed subscribes to active asset context with parsed messages
func (w *WebsocketClient) SubscribeToActiveAssetCtxParsed(coin string, callback func(WsGeneralActiveAssetCtx)) (int, error) {
	sub := NewParsedSubscription("activeAssetCtx", callback).WithCoin(coin)
	return w.SubscribeParsed(sub)
}

// matchesParsedSubscription determines if parsed data matches a subscription based on filters
func (w *WebsocketClient) matchesParsedSubscription(key subKey, parsedData any) bool {
	// If no filters are set, always match
	if key.coin == "" && key.user == "" && key.interval == "" {
		return true
	}

	// Handle different data types based on subscription type
	switch key.typ {
	case "trades":
		// Trades come as an array, check if any trade matches the coin filter
		if trades, ok := parsedData.([]WsTrade); ok && key.coin != "" {
			for _, trade := range trades {
				if trade.Coin == key.coin {
					return true
				}
			}
			return false
		}
	case "l2Book":
		// L2Book has a coin field
		if book, ok := parsedData.(WsBook); ok && key.coin != "" {
			return book.Coin == key.coin
		}
	case "candle":
		// Candle has a symbol field that corresponds to coin
		if candles, ok := parsedData.(WsCandle); ok {
			if key.coin != "" && candles.S != key.coin {
				return false
			}
			if key.interval != "" && candles.I != key.interval {
				return false
			}
			return true
		}
	case "userFills":
		// UserFills are already filtered by user at the subscription level
		// The websocket only sends fills for the specific user, so no additional filtering needed
		return true
	case "orderUpdates":
		// OrderUpdates come as an array, check if any order matches the coin filter
		if orders, ok := parsedData.([]WsOrder); ok && key.coin != "" {
			for _, order := range orders {
				if order.Order.Coin == key.coin {
					return true
				}
			}
			return false
		}
	case "activeAssetCtx":
		if activeAssetCtx, ok := parsedData.(WsGeneralActiveAssetCtx); ok && key.coin != "" {
			return activeAssetCtx.Coin == key.coin
		}
	}

	// For other types or if no specific filtering is needed, return true
	return true
}
