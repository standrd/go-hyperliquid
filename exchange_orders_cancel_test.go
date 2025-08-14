package hyperliquid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var recordForDebug = false

func TestCancel(t *testing.T) {
	type tc struct {
		name         string
		cassetteName string
		// If placeFirst is true, we first place a resting order and use its OID.
		placeFirst bool
		order      CreateOrderRequest
		coin       string
		oid        int64 // used only when placeFirst == false
		// If doubleCancel is true, we attempt to cancel the same OID twice to exercise the error path.
		doubleCancel bool
		wantErr      string
		record       bool
	}

	cases := []tc{
		{
			name:         "cancel resting order by oid",
			cassetteName: "Cancel",
			placeFirst:   true,
			order: CreateOrderRequest{
				Coin:  "DOGE",
				IsBuy: true,
				Size:  45,
				Price: 0.12330, // low so it stays resting
				OrderType: OrderType{
					Limit: &LimitOrderType{Tif: TifGtc},
				},
			},
			coin:   "DOGE",
			record: recordForDebug,
		},
		{
			name:         "double cancel returns error on second attempt",
			cassetteName: "Cancel",
			placeFirst:   true,
			order: CreateOrderRequest{
				Coin:  "DOGE",
				IsBuy: true,
				Size:  45,
				Price: 0.12330,
				OrderType: OrderType{
					Limit: &LimitOrderType{Tif: TifGtc},
				},
			},
			coin:         "DOGE",
			doubleCancel: true,
			wantErr:      "already canceled",
			record:       recordForDebug,
		},
		{
			name:         "cancel non-existent oid",
			cassetteName: "Cancel",
			placeFirst:   false,
			coin:         "DOGE",
			oid:          1,
			wantErr:      "Order was never placed, already canceled, or filled.",
			record:       recordForDebug,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			initRecorder(tt, tc.record, tc.cassetteName)

			exchange, err := newExchange(
				"0x38d55ff1195c57b9dbc8a72c93119500f1fcd47a33f98149faa18d2fc37932fa",
				TestnetAPIURL)
			require.NoError(t, err)

			oid := tc.oid
			if tc.placeFirst {
				placed, err := exchange.Order(tc.order, nil)
				require.NoError(tt, err)
				require.NotNil(tt, placed.Resting, "expected resting order so it can be canceled")
				oid = placed.Resting.Oid
			}

			// First cancel
			resp, err := exchange.Cancel(tc.coin, oid)
			if tc.wantErr != "" && !tc.doubleCancel {
				require.Error(tt, err)
				require.Contains(tt, err.Error(), tc.wantErr)
				return
			}
			require.NoError(tt, err)
			tt.Logf("cancel response: %+v", resp)

			// Optional second cancel to test error path
			if tc.doubleCancel {
				resp2, err2 := exchange.Cancel(tc.coin, oid)
				require.Error(tt, err2, "expected error on second cancel")
				if tc.wantErr != "" {
					require.Contains(tt, err2.Error(), tc.wantErr)
				}
				tt.Logf("second cancel response: %+v, err: %v", resp2, err2)
			}
		})
	}
}
