package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCollectDueSnapshots_OnlySamplesInsideTolerance(t *testing.T) {
	now := time.Date(2026, 7, 9, 2, 0, 0, 0, time.UTC)
	trades := []whaleTrade{
		{
			Wallet:  "0x1111111111111111111111111111111111111111",
			Side:    "BUY",
			AssetID: "asset-1",
			TradeID: "tx-1",
			Price:   0.40,
			Size:    1000,
			Time:    now.Add(-5*time.Minute - 30*time.Second),
		},
		{
			Wallet:  "0x2222222222222222222222222222222222222222",
			Side:    "BUY",
			AssetID: "asset-2",
			TradeID: "tx-2",
			Price:   0.50,
			Size:    1000,
			Time:    now.Add(-8 * time.Minute),
		},
	}

	got, err := collectDueSnapshots(context.Background(), trades, nil, []time.Duration{5 * time.Minute}, 2*time.Minute, now, func(context.Context, string) (float64, error) {
		return 0.46, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("snapshots len=%d, want 1", len(got))
	}
	if got[0].SignalID != signalID(trades[0]) {
		t.Fatalf("snapshot signal=%s, want first trade", got[0].SignalID)
	}
	if got[0].DeltaPP != 6 {
		t.Fatalf("DeltaPP=%.2f, want 6.00", got[0].DeltaPP)
	}
}

func TestCollectDueSnapshots_SkipsExistingSnapshot(t *testing.T) {
	now := time.Date(2026, 7, 9, 2, 0, 0, 0, time.UTC)
	trade := whaleTrade{
		Wallet:  "0x1111111111111111111111111111111111111111",
		AssetID: "asset-1",
		TradeID: "tx-1",
		Price:   0.40,
		Size:    1000,
		Time:    now.Add(-5 * time.Minute),
	}
	existing := map[string]struct{}{snapshotKey(signalID(trade), 5*time.Minute): {}}

	got, err := collectDueSnapshots(context.Background(), []whaleTrade{trade}, existing, []time.Duration{5 * time.Minute}, 2*time.Minute, now, func(context.Context, string) (float64, error) {
		t.Fatal("fetch should not be called for existing snapshot")
		return 0, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("snapshots len=%d, want 0", len(got))
	}
}

func TestApplySnapshotMetas_UsesCurrentWalletList(t *testing.T) {
	snaps := []edgeSnapshot{{
		Wallet: "0x1111111111111111111111111111111111111111",
		List:   "core",
		Tier:   "A",
	}}
	metas := map[string]walletMeta{
		"0x1111111111111111111111111111111111111111": {List: "sports", Tier: "B"},
	}

	applySnapshotMetas(snaps, metas)
	if snaps[0].List != "sports" || snaps[0].Tier != "B" {
		t.Fatalf("snapshot meta=%s/%s, want sports/B", snaps[0].List, snaps[0].Tier)
	}
}

func TestLoadWalletMetas_CommaSeparatedFiles(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "push.txt")
	second := filepath.Join(dir, "observe.txt")
	if err := os.WriteFile(first, []byte("0x1111111111111111111111111111111111111111 # list=core tier=A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("0x2222222222222222222222222222222222222222 # list=tape_observe tier=TAPE\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadWalletMetas(first + "," + second)
	if err != nil {
		t.Fatal(err)
	}
	if got["0x1111111111111111111111111111111111111111"].List != "core" {
		t.Fatalf("first list=%q, want core", got["0x1111111111111111111111111111111111111111"].List)
	}
	if got["0x2222222222222222222222222222222222222222"].List != "tape_observe" {
		t.Fatalf("second list=%q, want tape_observe", got["0x2222222222222222222222222222222222222222"].List)
	}
}

func TestLoadTapeTrades_UsesAllowlistAndObservationMeta(t *testing.T) {
	dir := t.TempDir()
	tape := filepath.Join(dir, "sports_tape.jsonl")
	data := "" +
		`{"time":"2026-07-09T02:00:00Z","wallet":"0x1111111111111111111111111111111111111111","side":"BUY","notional":6000,"price":0.42,"outcome":"Yes","market":"Soccer","asset":"asset-1","transaction":"tx-1"}` + "\n" +
		`{"time":"2026-07-09T02:00:00Z","wallet":"0x3333333333333333333333333333333333333333","side":"BUY","notional":9000,"price":0.50,"outcome":"No","market":"Soccer","asset":"asset-2","transaction":"tx-2"}` + "\n" +
		`{"time":"2026-07-09T02:00:00Z","wallet":"0x2222222222222222222222222222222222222222","side":"SELL","notional":7000,"price":0.51,"outcome":"No","market":"Soccer","asset":"asset-3","transaction":"tx-3"}` + "\n"
	if err := os.WriteFile(tape, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	metas := map[string]walletMeta{
		"0x1111111111111111111111111111111111111111": {List: "tape_observe", Tier: "TAPE"},
		"0x2222222222222222222222222222222222222222": {List: "tape_observe", Tier: "TAPE"},
	}

	got, err := loadTapeTrades(tape, metas, 5000)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("tape trades len=%d, want 1", len(got))
	}
	if got[0].Wallet != "0x1111111111111111111111111111111111111111" || got[0].Action != "tape" {
		t.Fatalf("trade=%+v, want allowlisted tape BUY", got[0])
	}
	if got[0].List != "tape_observe" || got[0].Tier != "TAPE" {
		t.Fatalf("meta=%s/%s, want tape_observe/TAPE", got[0].List, got[0].Tier)
	}
}
