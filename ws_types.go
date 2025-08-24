package hyperliquid

import (
	"encoding/json"
)

//go:generate easyjson

const (
	ChannelTrades  string = "trades"
	ChannelL2Book  string = "l2Book"
	ChannelCandle  string = "candle"
	ChannelAllMids string = "allMids"
)

//easyjson:json
type wsMessage struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

type Subscription struct {
	ID      string
	Payload any
	Close   func()
}

//easyjson:json
type WsCommand struct {
	Method       string `json:"method"`
	Subscription any    `json:"subscription,omitempty"`
}

type subscriptable interface {
	Key() string
}

type (
	Trades  []Trade
	Candles []Candle
)

func (t Trades) Key() string {
	if len(t) == 0 {
		return ""
	}
	return keyTrades(t[0].Coin)
}

func (c Candles) Key() string {
	if len(c) == 0 {
		return ""
	}
	return key(ChannelCandle, c[0].Symbol, c[0].Interval)
}

func (c L2Book) Key() string {
	return key(ChannelL2Book, c.Coin)
}
