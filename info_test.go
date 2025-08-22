package hyperliquid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetaAndAssetCtxs(t *testing.T) {
	info := NewInfo(MainnetAPIURL, true, nil, nil)

	initRecorder(t, false, "MetaAndAssetCtxs")

	res, err := info.MetaAndAssetCtxs()
	t.Logf("res: %+v", res)
	t.Logf("err: %v", err)

	require.NoError(t, err)

	// Verify the response structure
	require.NotNil(t, res)
	require.NotNil(t, res.Meta.Universe)
	require.NotNil(t, res.Meta.MarginTables)
	require.NotNil(t, res.Ctxs)

	// Verify we have at least one asset in universe
	require.Greater(t, len(res.Meta.Universe), 0)
	require.NotEmpty(t, res.Meta.Universe[0].Name)

	// Test specific known assets from the cassette data
	var btcFound, ethFound bool
	for _, asset := range res.Meta.Universe {
		if asset.Name == "BTC" {
			btcFound = true
			require.Equal(t, 5, asset.SzDecimals)
		}
		if asset.Name == "ETH" {
			ethFound = true
			require.Equal(t, 4, asset.SzDecimals)
		}
	}
	require.True(t, btcFound, "BTC asset should be present in universe")
	require.True(t, ethFound, "ETH asset should be present in universe")

	// Verify we have at least one margin table
	require.Greater(t, len(res.Meta.MarginTables), 0)
	require.GreaterOrEqual(t, res.Meta.MarginTables[0].ID, 0)

	// Verify we have at least one margin tier
	require.Greater(t, len(res.Meta.MarginTables[0].MarginTiers), 0)

	// Test specific margin table structure
	for _, marginTable := range res.Meta.MarginTables {
		require.NotNil(t, marginTable)
		require.Greater(t, len(marginTable.MarginTiers), 0)
		for _, tier := range marginTable.MarginTiers {
			require.NotEmpty(t, tier.LowerBound)
			require.Greater(t, tier.MaxLeverage, 0)
		}
	}

	// Verify we have at least one context
	require.Greater(t, len(res.Ctxs), 0)
	require.NotEmpty(t, res.Ctxs[0].MarkPx)
}

func TestSpotMetaAndAssetCtxs(t *testing.T) {
	info := NewInfo(MainnetAPIURL, true, nil, nil)

	initRecorder(t, false, "SpotMetaAndAssetCtxs")

	res, err := info.SpotMetaAndAssetCtxs()
	t.Logf("res: %+v", res)
	t.Logf("err: %v", err)

	require.NoError(t, err)

	// Verify the response structure
	require.NotNil(t, res)
	require.NotNil(t, res.Meta.Universe)
	require.NotNil(t, res.Meta.Tokens)
	require.NotNil(t, res.Ctxs)

	// Verify we have at least one asset in universe
	require.Greater(t, len(res.Meta.Universe), 0)
	require.NotEmpty(t, res.Meta.Universe[0].Name)

	// Test specific known assets from the cassette data
	var purrFound bool
	for _, asset := range res.Meta.Universe {
		if asset.Name == "PURR/USDC" {
			purrFound = true
			require.Equal(t, 0, asset.Index)
			require.True(t, asset.IsCanonical)
			require.Equal(t, []int{1, 0}, asset.Tokens)
		}
	}
	require.True(t, purrFound, "PURR/USDC asset should be present in universe")

	// Verify we have at least one token
	require.Greater(t, len(res.Meta.Tokens), 0)
	require.NotEmpty(t, res.Meta.Tokens[0].Name)

	// Test specific known tokens from the cassette data
	var usdcFound, purrTokenFound bool
	for _, token := range res.Meta.Tokens {
		if token.Name == "USDC" {
			usdcFound = true
			require.Equal(t, 8, token.SzDecimals)
			require.Equal(t, 8, token.WeiDecimals)
			require.Equal(t, 0, token.Index)
			require.True(t, token.IsCanonical)
		}
		if token.Name == "PURR" {
			purrTokenFound = true
			require.Equal(t, 0, token.SzDecimals)
			require.Equal(t, 5, token.WeiDecimals)
			require.Equal(t, 1, token.Index)
			require.True(t, token.IsCanonical)
		}
	}
	require.True(t, usdcFound, "USDC token should be present in tokens")
	require.True(t, purrTokenFound, "PURR token should be present in tokens")

	// Verify we have at least one context
	require.Greater(t, len(res.Ctxs), 0)
	require.NotEmpty(t, res.Ctxs[0].Coin)
}

func TestMeta(t *testing.T) {
	info := NewInfo(MainnetAPIURL, true, nil, nil)

	initRecorder(t, false, "Meta")

	res, err := info.Meta()
	t.Logf("res: %+v", res)
	t.Logf("err: %v", err)

	require.NoError(t, err)

	// Verify the response structure
	require.NotNil(t, res)
	require.NotNil(t, res.Universe)
	require.NotNil(t, res.MarginTables)

	// Verify we have at least one asset in universe
	require.Greater(t, len(res.Universe), 0)
	require.NotEmpty(t, res.Universe[0].Name)

	// Test specific known assets from the cassette data
	var btcFound, ethFound bool
	for _, asset := range res.Universe {
		if asset.Name == "BTC" {
			btcFound = true
			require.Equal(t, 5, asset.SzDecimals)
		}
		if asset.Name == "ETH" {
			ethFound = true
			require.Equal(t, 4, asset.SzDecimals)
		}
	}
	require.True(t, btcFound, "BTC asset should be present in universe")
	require.True(t, ethFound, "ETH asset should be present in universe")

	// Verify we have at least one margin table
	require.Greater(t, len(res.MarginTables), 0)
	require.GreaterOrEqual(t, res.MarginTables[0].ID, 0)

	// Test specific margin table structure
	for _, marginTable := range res.MarginTables {
		require.NotNil(t, marginTable)
		require.Greater(t, len(marginTable.MarginTiers), 0)
		for _, tier := range marginTable.MarginTiers {
			require.NotEmpty(t, tier.LowerBound)
			require.Greater(t, tier.MaxLeverage, 0)
		}
	}
}

func TestSpotMeta(t *testing.T) {
	info := NewInfo(MainnetAPIURL, true, nil, nil)

	initRecorder(t, false, "SpotMeta")

	res, err := info.SpotMeta()
	t.Logf("res: %+v", res)
	t.Logf("err: %v", err)

	require.NoError(t, err)

	// Verify the response structure
	require.NotNil(t, res)
	require.NotNil(t, res.Universe)
	require.NotNil(t, res.Tokens)

	// Verify we have at least one asset in universe
	require.Greater(t, len(res.Universe), 0)
	require.NotEmpty(t, res.Universe[0].Name)

	// Test specific known assets from the cassette data
	var purrFound bool
	for _, asset := range res.Universe {
		if asset.Name == "PURR/USDC" {
			purrFound = true
			require.Equal(t, 0, asset.Index)
			require.True(t, asset.IsCanonical)
			require.Equal(t, []int{1, 0}, asset.Tokens)
		}
	}
	require.True(t, purrFound, "PURR/USDC asset should be present in universe")

	// Verify we have at least one token
	require.Greater(t, len(res.Tokens), 0)
	require.NotEmpty(t, res.Tokens[0].Name)

	// Test specific known tokens from the cassette data
	var usdcFound, purrTokenFound bool
	for _, token := range res.Tokens {
		if token.Name == "USDC" {
			usdcFound = true
			require.Equal(t, 8, token.SzDecimals)
			require.Equal(t, 8, token.WeiDecimals)
			require.Equal(t, 0, token.Index)
			require.True(t, token.IsCanonical)
		}
		if token.Name == "PURR" {
			purrTokenFound = true
			require.Equal(t, 0, token.SzDecimals)
			require.Equal(t, 5, token.WeiDecimals)
			require.Equal(t, 1, token.Index)
			require.True(t, token.IsCanonical)
		}
	}
	require.True(t, usdcFound, "USDC token should be present in tokens")
	require.True(t, purrTokenFound, "PURR token should be present in tokens")
}

func TestQueryOrderByOid(t *testing.T) {
	type tc struct {
		name         string
		cassetteName string
		user         string
		oid          int64
		expected     *OrderQueryResult
		wantErr      string
		record       bool
		useTestnet   bool
	}

	info := NewInfo(MainnetAPIURL, true, nil, nil)

	cases := []tc{
		{
			name:         "TRX unknown order",
			cassetteName: "QueryOrderByOid",
			user:         "0x31ca8395cf837de08b24da3f660e77761dfb974b",
			oid:          141622259364,
			expected: &OrderQueryResult{
				Status: OrderQueryStatusError,
				Order:  OrderQueryResponse{},
			},
			record: false,
		},
		{
			name:         "SAND unknown order",
			cassetteName: "QueryOrderByOid",
			user:         "0x31ca8395cf837de08b24da3f660e77761dfb974b",
			oid:          141623226620,
			expected: &OrderQueryResult{
				Status: OrderQueryStatusError,
				Order:  OrderQueryResponse{},
			},
			record: false,
		},
		{
			name:         "User 0x8e0C473fed9630906779f982Cd0F80Cb7011812D order 37907159219",
			cassetteName: "QueryOrderByOid",
			user:         "0x8e0C473fed9630906779f982Cd0F80Cb7011812D",
			oid:          37907159219,
			expected: &OrderQueryResult{
				Status: OrderQueryStatusSuccess,
				Order: OrderQueryResponse{
					Order: QueriedOrder{
						Coin:             "ETH",
						Side:             OrderSideBid,
						LimitPx:          "4650.4",
						Sz:               "0.0",
						Oid:              37907159219,
						Timestamp:        1755857898644,
						TriggerCondition: "N/A",
						IsTrigger:        false,
						TriggerPx:        "0.0",
						IsPositionTpsl:   false,
						ReduceOnly:       false,
						OrderType:        "Market",
						OrigSz:           "0.0025",
						Tif:              "FrontendMarket",
						Cloid:            nil,
					},
					Status:          OrderStatusValueFilled,
					StatusTimestamp: 1755857898644,
				},
			},
			record:     false,
			useTestnet: true,
		},
		{
			name:         "User 0x8e0C473fed9630906779f982Cd0F80Cb7011812D order 37907165748",
			cassetteName: "QueryOrderByOid",
			user:         "0x8e0C473fed9630906779f982Cd0F80Cb7011812D",
			oid:          37907165748,
			expected: &OrderQueryResult{
				Status: OrderQueryStatusSuccess,
				Order: OrderQueryResponse{
					Order: QueriedOrder{
						Coin:             "ETH",
						Side:             OrderSideAsk,
						LimitPx:          "3960.7",
						Sz:               "0.0",
						Oid:              37907165748,
						Timestamp:        1755857910772,
						TriggerCondition: "N/A",
						IsTrigger:        false,
						TriggerPx:        "0.0",
						IsPositionTpsl:   false,
						ReduceOnly:       true,
						OrderType:        "Market",
						OrigSz:           "0.0025",
						Tif:              "FrontendMarket",
						Cloid:            nil,
					},
					Status:          OrderStatusValueFilled,
					StatusTimestamp: 1755857910772,
				},
			},
			record:     false,
			useTestnet: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			initRecorder(tt, tc.record, tc.cassetteName)

			var infoInstance *Info
			if tc.useTestnet {
				infoInstance = NewInfo(TestnetAPIURL, true, nil, nil)
			} else {
				infoInstance = info
			}

			res, err := infoInstance.QueryOrderByOid(tc.user, tc.oid)
			tt.Logf("res: %+v", res)
			tt.Logf("err: %v", err)

			if tc.wantErr != "" {
				require.Error(tt, err)
				require.Contains(tt, err.Error(), tc.wantErr)
				return
			} else {
				require.NoError(tt, err)
			}

			if err == nil {
				require.NotNil(tt, res)
				require.Equal(tt, tc.expected.Status, res.Status)

				// If order is found, compare order details
				if res.Status == OrderQueryStatusSuccess {
					require.Equal(tt, tc.expected.Order.Status, res.Order.Status)
					require.Equal(tt, tc.expected.Order.StatusTimestamp, res.Order.StatusTimestamp)

					// Compare order details
					expectedOrder := tc.expected.Order.Order
					actualOrder := res.Order.Order
					require.Equal(tt, expectedOrder.Coin, actualOrder.Coin)
					require.Equal(tt, expectedOrder.Side, actualOrder.Side)
					require.Equal(tt, expectedOrder.LimitPx, actualOrder.LimitPx)
					require.Equal(tt, expectedOrder.Sz, actualOrder.Sz)
					require.Equal(tt, expectedOrder.Oid, actualOrder.Oid)
					require.Equal(tt, expectedOrder.Timestamp, actualOrder.Timestamp)
					require.Equal(tt, expectedOrder.TriggerCondition, actualOrder.TriggerCondition)
					require.Equal(tt, expectedOrder.IsTrigger, actualOrder.IsTrigger)
					require.Equal(tt, expectedOrder.TriggerPx, actualOrder.TriggerPx)
					require.Equal(tt, expectedOrder.IsPositionTpsl, actualOrder.IsPositionTpsl)
					require.Equal(tt, expectedOrder.ReduceOnly, actualOrder.ReduceOnly)
					require.Equal(tt, expectedOrder.OrderType, actualOrder.OrderType)
					require.Equal(tt, expectedOrder.OrigSz, actualOrder.OrigSz)
					require.Equal(tt, expectedOrder.Tif, actualOrder.Tif)
					require.Equal(tt, expectedOrder.Cloid, actualOrder.Cloid)
				}
			}
		})
	}
}
