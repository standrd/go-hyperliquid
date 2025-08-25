package hyperliquid

import "github.com/sonirico/vago/fp"

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
	return keyCandles(c[0].Symbol, c[0].Interval)
}

func (c L2Book) Key() string {
	return keyL2Book(c.Coin)
}

func (a AllMids) Key() string {
	return keyAllMids(fp.None[string]())
}
