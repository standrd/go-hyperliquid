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
