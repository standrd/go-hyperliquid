package examples

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sonirico/go-hyperliquid"
)

// TradeHandler implements WSMessageHandler for Trade messages
type TradeHandler struct {
	trades []hyperliquid.WsTrade
	mu     sync.Mutex
}

func (h *TradeHandler) Handle(trades []hyperliquid.WsTrade) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.trades = append(h.trades, trades...)
}

func (h *TradeHandler) GetTrades() []hyperliquid.WsTrade {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.trades
}

// OrderbookHandler implements WSMessageHandler for L2Book messages
type OrderbookHandler struct {
	orderbooks []hyperliquid.WsBook
	mu         sync.Mutex
}

func (h *OrderbookHandler) Handle(orderbook hyperliquid.WsBook) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.orderbooks = append(h.orderbooks, orderbook)
}

func (h *OrderbookHandler) GetOrderbooks() []hyperliquid.WsBook {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.orderbooks
}

func TestParsedWebsocket(t *testing.T) {
	ws := hyperliquid.NewWebsocketClient(hyperliquid.MainnetAPIURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("Connecting to websocket")
	if err := ws.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer ws.Close()

	t.Log("Connected to websocket")

	// Test parsed trades subscription
	t.Run("parsed trades subscription", func(t *testing.T) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		received := make(chan struct{})

		tradeHandler := &TradeHandler{}

		_, err := ws.SubscribeToTradesParsed("BTC", tradeHandler.Handle)
		if err != nil {
			t.Fatalf("Failed to subscribe to trades: %v", err)
		}

		go func() {
			// Wait for at least one trade message
			time.Sleep(5 * time.Second)
			wg.Done()
		}()

		go func() {
			wg.Wait()
			close(received)
		}()

		select {
		case <-received:
			trades := tradeHandler.GetTrades()
			if len(trades) > 0 {
				t.Logf("Received %d trades", len(trades))
				for i, trade := range trades {
					if i >= 3 { // Log only first 3 trades
						break
					}
					t.Logf("Trade %d: Coin=%s, Side=%s, Price=%s, Size=%s",
						i+1, trade.Coin, trade.Side, trade.Px, trade.Sz)
				}
			} else {
				t.Log("No trades received yet")
			}
		case <-ctx.Done():
			t.Fatal("Timeout waiting for trades messages")
		}
	})

	// Test parsed orderbook subscription
	t.Run("parsed orderbook subscription", func(t *testing.T) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		received := make(chan struct{})

		orderbookHandler := &OrderbookHandler{}

		_, err := ws.SubscribeToOrderbookParsed("BTC", orderbookHandler.Handle)
		if err != nil {
			t.Fatalf("Failed to subscribe to orderbook: %v", err)
		}

		go func() {
			// Wait for at least one orderbook message
			time.Sleep(5 * time.Second)
			wg.Done()
		}()

		go func() {
			wg.Wait()
			close(received)
		}()

		select {
		case <-received:
			orderbooks := orderbookHandler.GetOrderbooks()
			if len(orderbooks) > 0 {
				t.Logf("Received %d orderbook updates", len(orderbooks))
				for i, orderbook := range orderbooks {
					if i >= 2 { // Log only first 2 orderbooks
						break
					}
					t.Logf("Orderbook %d: Coin=%s, Levels=%d",
						i+1, orderbook.Coin, len(orderbook.Levels))
				}
			} else {
				t.Log("No orderbook updates received yet")
			}
		case <-ctx.Done():
			t.Fatal("Timeout waiting for orderbook messages")
		}
	})
}

// Example of using the parsed subscription with a custom handler
func ExampleParsedSubscription() {
	ws := hyperliquid.NewWebsocketClient(hyperliquid.MainnetAPIURL)

	// Create a custom handler for trades
	tradeHandler := &TradeHandler{}

	// Subscribe with automatic parsing - the compiler ensures we're handling the correct type
	_, err := ws.SubscribeToTradesParsed("BTC", tradeHandler.Handle)
	if err != nil {
		panic(err)
	}

	// The handler will receive []Trade messages, not raw WSMessage
	// This provides compile-time type safety and automatic JSON parsing
}
