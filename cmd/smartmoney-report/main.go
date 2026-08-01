package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/15529214579/polymarket-go/internal/journal"
	"github.com/15529214579/polymarket-go/internal/paperreport"
	"github.com/15529214579/polymarket-go/internal/strategy"
)

func main() {
	journalDir := flag.String("journal", "db/smartmoney-paper/journal", "paper journal directory")
	positionsPath := flag.String("positions", "db/smartmoney-paper/positions.json", "paper positions state")
	markdownOut := flag.String("out", "reports/smartmoney-paper-pnl.md", "markdown output; empty disables")
	jsonOut := flag.String("json_out", "db/smartmoney-paper/pnl-report.json", "JSON output; empty disables")
	flag.Parse()

	trades, err := journal.ReadAll(*journalDir)
	if err != nil {
		fatalf("read journal: %v", err)
	}
	open, err := readOpenPositions(*positionsPath)
	if err != nil {
		fatalf("read positions: %v", err)
	}
	report := paperreport.Analyze(trades, open)
	markdown := paperreport.FormatMarkdown(report)
	if *markdownOut != "" {
		if err := writeFile(*markdownOut, []byte(markdown)); err != nil {
			fatalf("write markdown: %v", err)
		}
	}
	if *jsonOut != "" {
		body, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fatalf("marshal JSON: %v", err)
		}
		if err := writeFile(*jsonOut, append(body, '\n')); err != nil {
			fatalf("write JSON: %v", err)
		}
	}
	fmt.Print(markdown)
}

func readOpenPositions(path string) ([]strategy.Position, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var state struct {
		Open []strategy.Position `json:"open"`
	}
	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return nil, err
	}
	return state.Open, nil
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
