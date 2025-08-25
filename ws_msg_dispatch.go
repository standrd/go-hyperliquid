package hyperliquid

import (
	"encoding/json"
	"fmt"
)

type msgDispatcher interface {
	Dispatch(subs []*uniqSubscriber, msg wsMessage) error
}

type msgDispatcherFunc[T any] func(subs []*uniqSubscriber, msg wsMessage) error

func (d msgDispatcherFunc[T]) Dispatch(subs []*uniqSubscriber, msg wsMessage) error {
	return d(subs, msg)
}

func NewMsgDispatcher[T subscriptable](channel string) msgDispatcher {
	return msgDispatcherFunc[T](func(subs []*uniqSubscriber, msg wsMessage) error {
		if msg.Channel != channel {
			return nil
		}

		var x T
		if err := json.Unmarshal(msg.Data, &x); err != nil {
			return fmt.Errorf("failed to unmarshal message: %v", err)
		}

		for _, subscriber := range subs {
			if subscriber.id == x.Key() {
				subscriber.dispatch(x)
			}
		}

		return nil
	})
}

func NewNoopDispatcher() msgDispatcher {
	return msgDispatcherFunc[any](func(subs []*uniqSubscriber, msg wsMessage) error {
		println(string(msg.Data))
		return nil
	})
}

func NewUserSpecificDispatcher[T any](channel string) msgDispatcher {
	return msgDispatcherFunc[T](func(subs []*uniqSubscriber, msg wsMessage) error {
		if msg.Channel != channel {
			return nil
		}

		var x T
		if err := json.Unmarshal(msg.Data, &x); err != nil {
			return fmt.Errorf("failed to unmarshal message: %v", err)
		}

		// For user-specific channels, dispatch to all subscribers since
		// the message doesn't contain user identification info
		for _, subscriber := range subs {
			subscriber.dispatch(x)
		}

		return nil
	})
}
