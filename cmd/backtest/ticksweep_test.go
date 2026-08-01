package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTickRows(t *testing.T, path string, rows []tickRow) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create ticks: %v", err)
	}
	enc := json.NewEncoder(f)
	for _, row := range rows {
		if err := enc.Encode(row); err != nil {
			_ = f.Close()
			t.Fatalf("encode tick: %v", err)
		}
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close ticks: %v", err)
	}
}

func TestLoadTickPathsRejectsMixedAssets(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC()
	writeTickRows(t, filepath.Join(dir, "clean.jsonl"), []tickRow{
		{PosID: "p10", AssetID: "asset-a", Time: now, Mid: 0.50},
		{PosID: "p10", AssetID: "asset-a", Time: now.Add(time.Second), Mid: 0.51},
	})
	writeTickRows(t, filepath.Join(dir, "polluted.jsonl"), []tickRow{
		{PosID: "p11", AssetID: "asset-a", Time: now, Mid: 0.50},
		{PosID: "p11", AssetID: "asset-b", Time: now.Add(time.Second), Mid: 0.51},
	})

	paths, err := loadTickPaths(dir)
	if err != nil {
		t.Fatalf("loadTickPaths: %v", err)
	}
	if len(paths) != 1 || paths[0].PosID != "p10" {
		t.Fatalf("paths=%+v, want only clean p10", paths)
	}
}

func TestReplayPathUsesElapsedTimeForTimeout(t *testing.T) {
	now := time.Now().UTC()
	p := posPath{
		EntryMid: 0.50,
		Ticks:    []float64{0.50, 0.51, 0.60},
		Times:    []time.Time{now, now.Add(5 * time.Minute), now.Add(10 * time.Minute)},
		Duration: 10 * time.Minute,
	}

	pnl, exit, held := replayPath(p, 0, 0, 10*60, 0)
	if exit != "timeout" || held != 600 {
		t.Fatalf("exit=%s held=%d, want timeout at 600s", exit, held)
	}
	if pnl < 0.99 || pnl > 1.01 {
		t.Fatalf("pnl=%f, want $1.00 on a $5 normalized position", pnl)
	}
}

func TestLoadJournalTradesIgnoresNonTradeJSONL(t *testing.T) {
	dir := t.TempDir()
	trade := []byte("{\"id\":\"p1\",\"pnl_usd\":1}\n")
	if err := os.WriteFile(filepath.Join(dir, "trades-2026-08-01.jsonl"), trade, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "whale_trades.jsonl"), trade, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := loadJournalTrades(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "p1" {
		t.Fatalf("journal trades=%+v", got)
	}
}
