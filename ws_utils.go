package hyperliquid

import (
	"strings"

	"github.com/sonirico/vago/fp"
)

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

func keyAllMids(dex fp.Option[string]) string {
	// Unfortunately, "dex" parameter is not returned neither in subscription ACK nor in the
	// allMids message, no we are rendered unable to distinguish between different DEXes from
	// subscriber's standpoint.
	//if dex.IsNone() {
	//	return key(ChannelAllMids)
	//}
	//return key(ChannelAllMids, dex.UnwrapUnsafe())
	return key(ChannelAllMids)
}
