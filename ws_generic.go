package hyperliquid

import (
	"encoding/json"
)

// WSMessageHandler is a generic interface for handling WebSocket messages
type WSMessageHandler[T any] interface {
	Handle(msg T)
}

// ParsedWSMessage represents a WebSocket message with parsed data
type ParsedWSMessage[T any] struct {
	Channel string `json:"channel"`
	Data    T      `json:"data"`
}

// ParsedSubscriptionInterface provides a non-generic interface for parsed subscriptions
type ParsedSubscriptionInterface interface {
	GetSubscription() Subscription
	GetCallback() any
	GetUnmarshaler() func([]byte) (any, error)
}

// ParsedSubscription represents a type-safe subscription with a specific message type
type ParsedSubscription[T any] struct {
	Subscription
	callback func(T)
}

// NewParsedSubscription creates a new parsed subscription
func NewParsedSubscription[T any](subType string, callback func(T)) *ParsedSubscription[T] {
	return &ParsedSubscription[T]{
		Subscription: Subscription{Type: subType},
		callback:     callback,
	}
}

// WithCoin adds a coin filter to the subscription
func (ps *ParsedSubscription[T]) WithCoin(coin string) *ParsedSubscription[T] {
	ps.Coin = coin
	return ps
}

// WithUser adds a user filter to the subscription
func (ps *ParsedSubscription[T]) WithUser(user string) *ParsedSubscription[T] {
	ps.User = user
	return ps
}

// WithInterval adds an interval filter to the subscription
func (ps *ParsedSubscription[T]) WithInterval(interval string) *ParsedSubscription[T] {
	ps.Interval = interval
	return ps
}

// GetSubscription returns the underlying Subscription
func (ps *ParsedSubscription[T]) GetSubscription() Subscription {
	return ps.Subscription
}

// GetHandler returns the message handler
func (ps *ParsedSubscription[T]) GetCallback() any {
	return ps.callback
}

// GetUnmarshaler returns a function that unmarshals JSON data to the correct type
func (ps *ParsedSubscription[T]) GetUnmarshaler() func([]byte) (any, error) {
	return func(data []byte) (any, error) {
		var result T
		err := json.Unmarshal(data, &result)
		return result, err
	}
}

// DefaultUnmarshaler provides default JSON unmarshaling for types
type DefaultUnmarshaler[T any] struct{}

func (d DefaultUnmarshaler[T]) Unmarshal(data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}

// parsedSubscriptionCallback represents a type-safe subscription callback
type parsedSubscriptionCallback[T any] struct {
	id        int
	callback  any
	unmarshal func([]byte) (any, error)
}
