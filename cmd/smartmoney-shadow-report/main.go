package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/15529214579/polymarket-go/internal/shadowreport"
)

func main() {
	logPath := flag.String("log", "logs/smartmoney-paper.log", "smartmoney JSON log")
	reportOut := flag.String("out", "reports/smartmoney-exit-shadow.md", "Markdown report output")
	jsonOut := flag.String("json_out", "db/smartmoney-paper/exit-shadow-report.json", "JSON report output")
	flag.Parse()

	observations, err := readObservations(*logPath)
	if err != nil {
		fatalf("read shadow observations: %v", err)
	}
	report := shadowreport.Analyze(observations)
	markdown := shadowreport.FormatMarkdown(report)
	if err := writeFile(*reportOut, []byte(markdown)); err != nil {
		fatalf("write Markdown: %v", err)
	}
	body, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fatalf("marshal JSON: %v", err)
	}
	if err := writeFile(*jsonOut, append(body, '\n')); err != nil {
		fatalf("write JSON: %v", err)
	}
	fmt.Printf("smartmoney-shadow-report samples=%d report=%s\n", report.Samples, *reportOut)
}

func readObservations(path string) ([]shadowreport.Observation, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var observations []shadowreport.Observation
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 2*1024*1024)
	for sc.Scan() {
		var row struct {
			Time          time.Time `json:"time"`
			Message       string    `json:"msg"`
			PosID         string    `json:"pos"`
			Policy        string    `json:"policy"`
			EntryTime     time.Time `json:"entry_time"`
			HeldSec       int       `json:"held_sec"`
			Question      string    `json:"question"`
			Source        string    `json:"source"`
			HoldProfile   string    `json:"hold_profile"`
			NetPnLUSD     float64   `json:"net_pnl_usd"`
			ActualCloseAt time.Time `json:"actual_close_at"`
		}
		if json.Unmarshal(sc.Bytes(), &row) != nil || row.Message != "copytrade_exit_shadow" {
			continue
		}
		entry := row.EntryTime
		if entry.IsZero() && !row.Time.IsZero() && row.HeldSec > 0 {
			entry = row.Time.Add(-time.Duration(row.HeldSec) * time.Second)
		}
		observations = append(observations, shadowreport.Observation{
			ObservedAt: row.Time, EntryTime: entry, ActualCloseAt: row.ActualCloseAt,
			PosID: row.PosID, Policy: row.Policy, Question: row.Question, Source: row.Source,
			HoldProfile: row.HoldProfile, NetPnLUSD: row.NetPnLUSD,
		})
	}
	return observations, sc.Err()
}

func writeFile(path string, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
