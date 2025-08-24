package hyperliquid

import "strings"

func key(args ...string) string {
	return strings.Join(args, ":")
}

func keyTrades(coin string) string {
	return key(ChannelTrades, coin)
}

func keyCandles(symbol, interval string) string {
	return key(ChannelCandle, symbol, interval)
}

func keyL2Book(coin string) string {
	return key(ChannelL2Book, coin)
}
