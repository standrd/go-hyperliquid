package examples

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/sonirico/go-hyperliquid"
)

func TestMetaAndAssetCtxs(t *testing.T) {
	godotenv.Overload()
	info := hyperliquid.NewInfo(hyperliquid.MainnetAPIURL, true, nil, nil)
	meta, err := info.MetaAndAssetCtxs()
	if err != nil {
		t.Fatalf("Failed to get meta: %v", err)
	}

	if meta.Meta.Universe == nil {
		t.Error("Expected non-nil universe")
	}

	if meta.Meta.MarginTables == nil {
		t.Error("Expected non-nil margin tables")
	}

	if meta.Ctxs == nil {
		t.Error("Expected non-nil contexts")
	}

	if len(meta.Meta.Universe) == 0 {
		t.Error("Expected at least one asset in universe")
	}

	if meta.Meta.Universe[0].Name == "" {
		t.Error("Expected name to be non-empty")
	}

	if len(meta.Meta.MarginTables) == 0 {
		t.Error("Expected at least one margin table")
	}

	if meta.Meta.MarginTables[0].ID < 0 {
		t.Error("Expected ID to be non-negative")
	}

	if len(meta.Meta.MarginTables[0].MarginTiers) == 0 {
		t.Error("Expected at least one margin tier")
	}

	if len(meta.Ctxs) == 0 {
		t.Error("Expected at least one context")
	}

	if meta.Ctxs[0].MarkPx == "" {
		t.Error("Expected mark price to be non-empty")
	}
}

func TestSpotMetaAndAssetCtxs(t *testing.T) {
	godotenv.Overload()

	info := hyperliquid.NewInfo(hyperliquid.MainnetAPIURL, true, nil, nil)
	spotMeta, err := info.SpotMetaAndAssetCtxs()
	if err != nil {
		t.Fatalf("Failed to get spot meta: %v", err)
	}

	if spotMeta.Meta.Universe == nil {
		t.Error("Expected non-nil universe")
	}

	if spotMeta.Meta.Tokens == nil {
		t.Error("Expected non-nil tokens")
	}

	if spotMeta.Ctxs == nil {
		t.Error("Expected non-nil contexts")
	}

	if len(spotMeta.Meta.Universe) == 0 {
		t.Error("Expected at least one asset in universe")
	}

	if spotMeta.Meta.Universe[0].Name == "" {
		t.Error("Expected name to be non-empty")
	}

	if len(spotMeta.Meta.Tokens) == 0 {
		t.Error("Expected at least one token")
	}

	if spotMeta.Meta.Tokens[0].Name == "" {
		t.Error("Expected name to be non-empty")
	}

	if len(spotMeta.Ctxs) == 0 {
		t.Error("Expected at least one context")
	}

	if spotMeta.Ctxs[0].Coin == "" {
		t.Error("Expected coin to be non-empty")
	}
}
