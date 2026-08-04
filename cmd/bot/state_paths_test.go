package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateTradingStatePaths(t *testing.T) {
	tests := []struct {
		name    string
		live    bool
		path    string
		wantErr bool
	}{
		{name: "paper scoped", path: "db/smartmoney-paper/positions.json"},
		{name: "live scoped", live: true, path: "db/live/positions.json"},
		{name: "live rejects paper", live: true, path: "db/smartmoney-paper/positions.json", wantErr: true},
		{name: "live rejects ambiguous legacy path", live: true, path: "db/positions.json", wantErr: true},
		{name: "paper rejects live", path: "db/smartmoney-live/positions.json", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTradingStatePaths(tt.live, tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateTradingStatePaths() err=%v wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTradingStatePathsRejectsLiveSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	liveLink := filepath.Join(dir, "live")
	if err := os.Symlink(target, liveLink); err != nil {
		t.Fatal(err)
	}
	if err := validateTradingStatePaths(true, filepath.Join(liveLink, "positions.json")); err == nil {
		t.Fatal("expected live symlink rejection")
	}
}
