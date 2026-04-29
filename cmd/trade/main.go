package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/order"
)

func main() {
	assetID := flag.String("asset", "", "ERC1155 token ID (YES side)")
	market := flag.String("market", "", "conditionID for logging")
	sizeUSD := flag.Float64("size", 0, "order size in USDC (0 = use all available)")
	limitPx := flag.Float64("price", 0, "limit price (0..1)")
	side := flag.String("side", "BUY", "BUY or SELL")
	negRisk := flag.Bool("negrisk", false, "use NegRisk exchange address")
	dryRun := flag.Bool("dry-run", false, "derive key + print intent, don't submit")
	autoWrap := flag.Bool("auto-wrap", true, "auto wrap USDC.e → pUSD if needed")
	rpcURL := flag.String("rpc", "", "Polygon RPC URL (default: polygon-rpc.com)")
	flag.Parse()

	if *assetID == "" || *limitPx <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: trade -asset <tokenID> -price <0..1> [-size <usd>] [-negrisk] [-dry-run]\n")
		os.Exit(1)
	}

	slog.Info("trade_init", "asset", *assetID, "size", *sizeUSD, "price", *limitPx, "side", *side, "negrisk", *negRisk)

	order.InitProxy()

	mnemonic, err := loadMnemonic()
	if err != nil {
		slog.Error("wallet_load_failed", "err", err)
		os.Exit(1)
	}
	wallet, err := order.NewWalletFromMnemonic(mnemonic, "")
	if err != nil {
		slog.Error("wallet_derive_failed", "err", err)
		os.Exit(1)
	}
	slog.Info("wallet_loaded", "address", wallet.Address().Hex())

	creds, err := order.DeriveAPIKey(order.ClobBaseURL, wallet)
	if err != nil {
		slog.Error("api_key_derive_failed", "err", err)
		os.Exit(1)
	}
	slog.Info("api_key_derived", "api_key", creds.APIKey)

	if *autoWrap {
		oc, err := order.NewOnChain(*rpcURL, wallet)
		if err != nil {
			slog.Error("onchain_init_failed", "err", err)
			os.Exit(1)
		}
		wrapCtx, wrapCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		if err := oc.WrapAll(wrapCtx); err != nil {
			slog.Error("auto_wrap_failed", "err", err)
			wrapCancel()
			os.Exit(1)
		}
		wrapCancel()

		approveCtx, approveCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		if err := oc.ApproveExchanges(approveCtx); err != nil {
			slog.Error("exchange_approve_failed", "err", err)
			approveCancel()
			os.Exit(1)
		}
		approveCancel()

		pusd, err := oc.PUSDBalance(context.Background())
		if err != nil {
			slog.Warn("pusd_balance_check_failed", "err", err)
		} else {
			pusdFloat, _ := new(big.Float).Quo(
				new(big.Float).SetInt(pusd),
				new(big.Float).SetFloat64(1e6),
			).Float64()
			slog.Info("pusd_balance", "raw", pusd.String(), "usd", fmt.Sprintf("%.2f", pusdFloat))

			if *sizeUSD <= 0 {
				*sizeUSD = pusdFloat * 0.98
				slog.Info("auto_size", "usd", *sizeUSD, "reason", "using 98% of pUSD balance")
			}
		}
	}

	if *sizeUSD <= 0 {
		fmt.Fprintf(os.Stderr, "ERROR: -size is required (or use -auto-wrap to auto-detect balance)\n")
		os.Exit(1)
	}

	intent := order.Intent{
		AssetID: *assetID,
		Market:  *market,
		Side:    order.Side(*side),
		SizeUSD: *sizeUSD,
		LimitPx: *limitPx,
		Type:    order.GTC,
	}

	fmt.Printf("\n=== ORDER INTENT ===\n")
	fmt.Printf("Asset:    %s\n", intent.AssetID)
	fmt.Printf("Side:     %s\n", intent.Side)
	fmt.Printf("Size:     $%.2f\n", intent.SizeUSD)
	fmt.Printf("Price:    %.4f\n", intent.LimitPx)
	fmt.Printf("Shares:   ~%.1f\n", intent.SizeUSD/intent.LimitPx)
	fmt.Printf("NegRisk:  %v\n", *negRisk)
	fmt.Printf("Exchange: %s\n", func() string {
		if *negRisk {
			return order.V2NegRiskExchangeAddress
		}
		return order.V2ExchangeAddress
	}())
	fmt.Printf("====================\n\n")

	if *dryRun {
		fmt.Println("DRY RUN — not submitting.")
		return
	}

	client := order.NewV2Client(wallet, creds, *negRisk)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	slog.Info("submitting_order")
	result, err := client.Submit(ctx, intent)
	if err != nil {
		slog.Error("order_failed", "err", err, "status", result.Status, "error_msg", result.Error)
		os.Exit(1)
	}

	fmt.Printf("\n=== ORDER RESULT ===\n")
	fmt.Printf("OrderID:  %s\n", result.OrderID)
	fmt.Printf("Status:   %s\n", result.Status)
	fmt.Printf("Filled:   %.4f shares\n", result.FilledSize)
	fmt.Printf("AvgPrice: %.4f\n", result.AvgPrice)
	fmt.Printf("FeeUSD:   $%.4f\n", result.FeeUSD)
	fmt.Printf("====================\n")

	if result.Status == order.StatusFilled {
		fmt.Println("\n✅ ORDER FILLED")
	} else {
		fmt.Printf("\n⚠️  Status: %s — %s\n", result.Status, result.Error)
	}
}

func loadMnemonic() (string, error) {
	out, err := exec.Command("bw", "get", "notes", "Polymarket-Go Wallet").Output()
	if err != nil {
		return "", fmt.Errorf("bw get notes: %w", err)
	}
	notes := string(out)
	lines := strings.Split(notes, "\n")
	for i, line := range lines {
		if strings.Contains(line, "助记词") && i+1 < len(lines) {
			mnemonic := strings.TrimSpace(lines[i+1])
			if mnemonic != "" {
				words := strings.Fields(mnemonic)
				if len(words) == 12 || len(words) == 24 {
					return mnemonic, nil
				}
			}
		}
	}
	return "", fmt.Errorf("mnemonic not found in Polymarket-Go Wallet notes")
}

func checkBalance(creds *order.APICredentials, addr string) {
	out, err := exec.Command("bw", "get", "item", "Polymarket-Go Wallet").Output()
	if err != nil {
		return
	}
	var item struct {
		Fields []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(out, &item); err != nil {
		return
	}
}
