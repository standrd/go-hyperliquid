package hyperliquid

import (
	"encoding/json"
)

//go:generate easyjson -all ws_types.go

type WSMessage struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

type Subscription struct {
	Type     string `json:"type"`
	Coin     string `json:"coin,omitempty"`
	User     string `json:"user,omitempty"`
	Interval string `json:"interval,omitempty"`
}

type subKey struct {
	typ      string
	coin     string
	user     string
	interval string
}

func (s Subscription) key() subKey {
	return subKey{
		typ:      s.Type,
		coin:     s.Coin,
		user:     s.User,
		interval: s.Interval,
	}
}

type WsCommand struct {
	Method       string        `json:"method"`
	Subscription *Subscription `json:"subscription,omitempty"`
}

type subscriptionCallback struct {
	id       int
	callback func(WSMessage)
}

type WsTrade struct {
	Coin  string   `json:"coin"`
	Side  string   `json:"side"`
	Px    string   `json:"px"`
	Sz    string   `json:"sz"`
	Hash  string   `json:"hash"`
	Time  int64    `json:"time"`
	Tid   int64    `json:"tid"`   // 50-bit hash of (buyer_oid, seller_oid)
	Users []string `json:"users"` // [buyer, seller]
}

type WsBook struct {
	Coin   string      `json:"coin"`
	Levels [][]WsLevel `json:"levels"`
	Time   int64       `json:"time"`
}

type WsBbo struct {
	Coin string     `json:"coin"`
	Time int64      `json:"time"`
	Bbo  []*WsLevel `json:"bbo"` // [bid, ask] - can be null
}

type WsLevel struct {
	Px string `json:"px"` // price
	Sz string `json:"sz"` // size
	N  int    `json:"n"`  // number of orders
}

type Notification struct {
	Notification string `json:"notification"`
}

type AllMids struct {
	Mids map[string]string `json:"mids"`
}
type WsCandle struct {
	T  int64  `json:"t"` // open millis
	TC int64  `json:"T"` // close millis
	S  string `json:"s"` // coin
	I  string `json:"i"` // interval
	O  string `json:"o"` // open price
	C  string `json:"c"` // close price
	H  string `json:"h"` // high price
	L  string `json:"l"` // low price
	V  string `json:"v"` // volume (base unit)
	N  int64  `json:"n"` // number of trades
}
type WsUserEvent struct {
	Fills         []WsFill          `json:"fills,omitempty"`
	Funding       *WsUserFunding    `json:"funding,omitempty"`
	Liquidation   *WsLiquidation    `json:"liquidation,omitempty"`
	NonUserCancel []WsNonUserCancel `json:"nonUserCancel,omitempty"`
}
type WsUserFills struct {
	IsSnapshot *bool    `json:"isSnapshot,omitempty"`
	User       string   `json:"user"`
	Fills      []WsFill `json:"fills"`
}
type WsFill struct {
	Coin          string           `json:"coin"`
	Px            string           `json:"px"` // price
	Sz            string           `json:"sz"` // size
	Side          string           `json:"side"`
	Time          int64            `json:"time"`
	StartPosition string           `json:"startPosition"`
	Dir           string           `json:"dir"` // used for frontend display
	ClosedPnl     string           `json:"closedPnl"`
	Hash          string           `json:"hash"`    // L1 transaction hash
	Oid           int64            `json:"oid"`     // order id
	Crossed       bool             `json:"crossed"` // whether order crossed the spread (was taker)
	Fee           string           `json:"fee"`     // negative means rebate
	Tid           int64            `json:"tid"`     // unique trade id
	Liquidation   *FillLiquidation `json:"liquidation,omitempty"`
	FeeToken      string           `json:"feeToken"`             // the token the fee was paid in
	BuilderFee    *string          `json:"builderFee,omitempty"` // amount paid to builder, also included in fee
}
type FillLiquidation struct {
	LiquidatedUser *string `json:"liquidatedUser,omitempty"`
	MarkPx         string  `json:"markPx"`
	Method         string  `json:"method"` // "market" | "backstop"
}
type WsUserFunding struct {
	Time        int64  `json:"time"`
	Coin        string `json:"coin"`
	Usdc        string `json:"usdc"`
	Szi         string `json:"szi"`
	FundingRate string `json:"fundingRate"`
}
type WsLiquidation struct {
	Lid                    int64  `json:"lid"`
	Liquidator             string `json:"liquidator"`
	LiquidatedUser         string `json:"liquidated_user"`
	LiquidatedNtlPos       string `json:"liquidated_ntl_pos"`
	LiquidatedAccountValue string `json:"liquidated_account_value"`
}
type WsNonUserCancel struct {
	Coin string `json:"coin"`
	Oid  int64  `json:"oid"`
}
type WsOrder struct {
	Order           WsBasicOrder `json:"order"`
	Status          string       `json:"status"` // See https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#query-order-status-by-oid-or-cloid for a list of possible values
	StatusTimestamp int64        `json:"statusTimestamp"`
}
type WsBasicOrder struct {
	Coin      string  `json:"coin"`
	Side      string  `json:"side"`
	LimitPx   string  `json:"limitPx"`
	Sz        string  `json:"sz"`
	Oid       int64   `json:"oid"`
	Timestamp int64   `json:"timestamp"`
	OrigSz    string  `json:"origSz"`
	Cloid     *string `json:"cloid,omitempty"`
}

type WsMergedActiveAssetCtx struct {
	WsSharedAssetCtx
	Funding           string `json:"funding,omitempty"`
	OpenInterest      string `json:"openInterest,omitempty"`
	OraclePx          string `json:"oraclePx,omitempty"`
	CirculatingSupply string `json:"circulatingSupply,omitempty"`
}

type WsGeneralActiveAssetCtx struct {
	Coin string                 `json:"coin"`
	Ctx  WsMergedActiveAssetCtx `json:"ctx"`
}

type WsActiveAssetCtx struct {
	Coin string          `json:"coin"`
	Ctx  WsPerpsAssetCtx `json:"ctx"`
}
type WsActiveSpotAssetCtx struct {
	Coin string         `json:"coin"`
	Ctx  WsSpotAssetCtx `json:"ctx"`
}
type WsSharedAssetCtx struct {
	DayNtlVlm string  `json:"dayNtlVlm"`
	PrevDayPx string  `json:"prevDayPx"`
	MarkPx    string  `json:"markPx"`
	MidPx     *string `json:"midPx,omitempty"`
}
type WsPerpsAssetCtx struct {
	WsSharedAssetCtx
	Funding      string `json:"funding"`
	OpenInterest string `json:"openInterest"`
	OraclePx     string `json:"oraclePx"`
}
type WsSpotAssetCtx struct {
	WsSharedAssetCtx
	CirculatingSupply string `json:"circulatingSupply"`
}
type WsActiveAssetData struct {
	User             string     `json:"user"`
	Coin             string     `json:"coin"`
	Leverage         WsLeverage `json:"leverage"`
	MaxTradeSzs      []string   `json:"maxTradeSzs"`
	AvailableToTrade []string   `json:"availableToTrade"`
}
type WsTwapSliceFill struct {
	Fill   WsFill `json:"fill"`
	TwapId int64  `json:"twapId"`
}
type WsUserTwapSliceFills struct {
	IsSnapshot     *bool             `json:"isSnapshot,omitempty"`
	User           string            `json:"user"`
	TwapSliceFills []WsTwapSliceFill `json:"twapSliceFills"`
}
type TwapState struct {
	Coin        string `json:"coin"`
	User        string `json:"user"`
	Side        string `json:"side"`
	Sz          string `json:"sz"`
	ExecutedSz  string `json:"executedSz"`
	ExecutedNtl string `json:"executedNtl"`
	Minutes     int    `json:"minutes"`
	ReduceOnly  bool   `json:"reduceOnly"`
	Randomize   bool   `json:"randomize"`
	Timestamp   int64  `json:"timestamp"`
}

type TwapStatus string

const (
	TwapStatusActivated  TwapStatus = "activated"
	TwapStatusTerminated TwapStatus = "terminated"
	TwapStatusFinished   TwapStatus = "finished"
	TwapStatusError      TwapStatus = "error"
)

type WsTwapHistory struct {
	State  TwapState `json:"state"`
	Status struct {
		Status      TwapStatus `json:"status"`
		Description string     `json:"description"`
	} `json:"status"`
	Time int64 `json:"time"`
}
type WsUserTwapHistory struct {
	IsSnapshot *bool           `json:"isSnapshot,omitempty"`
	User       string          `json:"user"`
	History    []WsTwapHistory `json:"history"`
}

type WebData2 struct {
	// TODO: not in docs?
}
type WsLeverage struct {
	// TODO: not in docs?
}
