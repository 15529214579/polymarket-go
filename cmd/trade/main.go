package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	nethttp "net/http"
	"os"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/order"
	"github.com/ethereum/go-ethereum/common"
)

const projectWalletAddress = "0x015282e9b720E072A9B87eEeaE738C6Bb039Bd9e"

func main() {
	assetID := flag.String("asset", "", "ERC1155 token ID (YES side)")
	market := flag.String("market", "", "conditionID for logging")
	sizeUSD := flag.Float64("size", 0, "order size in USDC (0 = use all available)")
	limitPx := flag.Float64("price", 0, "limit price (0..1)")
	side := flag.String("side", "BUY", "BUY or SELL")
	negRisk := flag.Bool("negrisk", false, "use NegRisk exchange address")
	dryRun := flag.Bool("dry-run", false, "validate and print intent locally; no network or wallet mutation")
	autoWrap := flag.Bool("auto-wrap", false, "explicitly wrap USDC.e and approve exchanges before submit")
	rpcURL := flag.String("rpc", "", "Polygon RPC URL (default: polygon-rpc.com)")
	hexKey := flag.String("key", "", "hex private key (bypasses Bitwarden mnemonic)")
	queryOpenOrders := flag.Bool("open-orders", false, "query CLOB open orders without cancelling")
	cancelOpenOrders := flag.Bool("cancel-open-orders", false, "explicitly cancel every open CLOB order")
	readiness := flag.Bool("readiness", false, "read-only wallet, L1/L2 auth, balance, allowance, and open-order checks")
	publicReadiness := flag.Bool("readiness-public", false, "read-only public CLOB and Polygon checks; no private key required")
	expectedAddress := flag.String("expected-address", projectWalletAddress, "wallet address that must match before authenticated operations")
	redeemAll := flag.Bool("redeem-all", false, "redeem all redeemable positions and exit")
	flag.Parse()

	managementMode := *queryOpenOrders || *cancelOpenOrders || *readiness || *publicReadiness || *redeemAll
	if !managementMode && (*assetID == "" || *limitPx <= 0) {
		fmt.Fprintf(os.Stderr, "Usage: trade -asset <tokenID> -price <0..1> [-size <usd>] [-negrisk] [-dry-run]\n")
		os.Exit(1)
	}
	if *expectedAddress != "" && !common.IsHexAddress(*expectedAddress) {
		fmt.Fprintf(os.Stderr, "ERROR: invalid -expected-address %q\n", *expectedAddress)
		os.Exit(1)
	}

	slog.Info("trade_init", "asset", *assetID, "size", *sizeUSD, "price", *limitPx, "side", *side, "negrisk", *negRisk)

	order.InitProxy()
	if *dryRun {
		intent, err := buildIntent(*assetID, *market, *side, *sizeUSD, *limitPx, *negRisk)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
		printIntent(intent)
		fmt.Println("DRY RUN - local validation only; no wallet access, API authentication, transaction, approval, cancellation, or order submission.")
		return
	}
	if *publicReadiness {
		if *expectedAddress == "" {
			fmt.Fprintln(os.Stderr, "ERROR: -expected-address is required for -readiness-public")
			os.Exit(1)
		}
		if err := runPublicReadiness(*rpcURL, common.HexToAddress(*expectedAddress)); err != nil {
			slog.Error("public_readiness_failed", "err", err)
			os.Exit(1)
		}
		return
	}

	var wallet *order.Wallet
	if *hexKey != "" {
		var err error
		wallet, err = order.NewWalletFromHexKey(*hexKey)
		if err != nil {
			slog.Error("wallet_from_key_failed", "err", err)
			os.Exit(1)
		}
	} else {
		mnemonic, err := order.LoadMnemonicFromBitwarden("Polymarket-Go Wallet", "mnemonic")
		if err != nil {
			slog.Error("wallet_load_failed", "err", err)
			os.Exit(1)
		}
		wallet, err = order.NewWalletFromMnemonic(mnemonic, "")
		if err != nil {
			slog.Error("wallet_derive_failed", "err", err)
			os.Exit(1)
		}
	}
	slog.Info("wallet_loaded", "address", wallet.Address().Hex())
	if *expectedAddress != "" && wallet.Address() != common.HexToAddress(*expectedAddress) {
		slog.Error("wallet_address_mismatch", "loaded", wallet.Address().Hex(), "expected", common.HexToAddress(*expectedAddress).Hex())
		os.Exit(1)
	}

	if *redeemAll {
		oc, err := order.NewOnChain(*rpcURL, wallet)
		if err != nil {
			slog.Error("onchain_init_failed", "err", err)
			os.Exit(1)
		}
		runRedeemAll(oc, wallet.Address().Hex())
		return
	}

	var creds *order.APICredentials
	var err error
	if *readiness || *queryOpenOrders || *cancelOpenOrders {
		creds, err = order.DeriveExistingAPIKey(order.ClobBaseURL, wallet)
	} else {
		creds, err = order.DeriveAPIKey(order.ClobBaseURL, wallet)
	}
	if err != nil {
		slog.Error("api_key_derive_failed", "err", err)
		os.Exit(1)
	}
	slog.Info("api_key_derived")
	client := order.NewV2Client(wallet, creds, *negRisk)

	if *readiness {
		if err := runAuthenticatedReadiness(*rpcURL, wallet, client); err != nil {
			slog.Error("authenticated_readiness_failed", "err", err)
			os.Exit(1)
		}
		return
	}

	if *queryOpenOrders {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		openOrders, err := client.GetOpenOrders(ctx)
		if err != nil {
			slog.Error("open_orders_query_failed", "err", err)
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(openOrders, "", "  ")
		fmt.Println(string(out))
		if !*cancelOpenOrders {
			return
		}
	}

	if *cancelOpenOrders {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := client.CancelAllOpen(ctx); err != nil {
			slog.Error("cancel_all_failed", "err", err)
			os.Exit(1)
		}
		fmt.Println("All open orders cancelled.")
		return
	}

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

	intent, err := buildIntent(*assetID, *market, *side, *sizeUSD, *limitPx, *negRisk)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
	printIntent(intent)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	refreshCtx, refreshCancel := context.WithTimeout(ctx, 15*time.Second)
	if err := client.RefreshBalanceAllowance(refreshCtx); err != nil {
		refreshCancel()
		slog.Error("balance_allowance_refresh_failed", "err", err)
		os.Exit(1)
	}
	refreshCancel()

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
		saveBuyTime(*assetID)
	} else {
		fmt.Printf("\n⚠️  Status: %s — %s\n", result.Status, result.Error)
	}
}

func buildIntent(assetID, market, side string, sizeUSD, limitPx float64, negRisk bool) (order.Intent, error) {
	side = strings.ToUpper(strings.TrimSpace(side))
	if assetID == "" {
		return order.Intent{}, fmt.Errorf("asset ID is required")
	}
	if side != string(order.Buy) && side != string(order.Sell) {
		return order.Intent{}, fmt.Errorf("side must be BUY or SELL")
	}
	if sizeUSD <= 0 {
		return order.Intent{}, fmt.Errorf("size must be greater than zero")
	}
	if limitPx <= 0 || limitPx >= 1 {
		return order.Intent{}, fmt.Errorf("price must be between zero and one")
	}
	return order.Intent{
		AssetID: assetID,
		Market:  market,
		Side:    order.Side(side),
		SizeUSD: sizeUSD,
		LimitPx: limitPx,
		Type:    order.GTC,
		NegRisk: negRisk,
	}, nil
}

func printIntent(intent order.Intent) {
	fmt.Printf("\n=== ORDER INTENT ===\n")
	fmt.Printf("Asset:    %s\n", intent.AssetID)
	fmt.Printf("Side:     %s\n", intent.Side)
	fmt.Printf("Size:     $%.2f\n", intent.SizeUSD)
	fmt.Printf("Price:    %.4f\n", intent.LimitPx)
	fmt.Printf("Shares:   ~%.1f\n", intent.SizeUSD/intent.LimitPx)
	fmt.Printf("NegRisk:  %v\n", intent.NegRisk)
	if intent.NegRisk {
		fmt.Printf("Exchange: %s\n", order.V2NegRiskExchangeAddress)
	} else {
		fmt.Printf("Exchange: %s\n", order.V2ExchangeAddress)
	}
	fmt.Printf("====================\n\n")
}

func runPublicReadiness(rpcURL string, address common.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	version, err := order.GetCLOBVersion(ctx, order.ClobBaseURL)
	if err != nil {
		return err
	}
	oc, err := order.NewReadOnlyOnChain(rpcURL, address)
	if err != nil {
		return err
	}
	snapshot, err := oc.ExchangeReadiness(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("CLOB version: %d (public endpoint OK)\n", version)
	fmt.Printf("Wallet: %s (public chain checks only)\n", address.Hex())
	printExchangeReadiness(snapshot)
	return nil
}

func runAuthenticatedReadiness(rpcURL string, wallet *order.Wallet, client *order.V2Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	version, err := order.GetCLOBVersion(ctx, order.ClobBaseURL)
	if err != nil {
		return err
	}
	openOrders, err := client.GetOpenOrders(ctx)
	if err != nil {
		return fmt.Errorf("L2 open-orders authentication: %w", err)
	}
	clobBalance, err := client.GetBalanceAllowance(ctx)
	if err != nil {
		return fmt.Errorf("L2 balance authentication: %w", err)
	}
	oc, err := order.NewReadOnlyOnChain(rpcURL, wallet.Address())
	if err != nil {
		return err
	}
	snapshot, err := oc.ExchangeReadiness(ctx)
	if err != nil {
		return err
	}
	count := openOrders.Count
	if count == 0 {
		count = len(openOrders.Data)
	}
	fmt.Printf("CLOB version: %d\n", version)
	fmt.Printf("Wallet: %s (expected address matched)\n", wallet.Address().Hex())
	fmt.Println("L1 existing API-key derivation: OK")
	fmt.Printf("L2 signed reads: OK (open orders: %d, cached collateral balance raw: %s)\n", count, clobBalance.Balance)
	printExchangeReadiness(snapshot)
	fmt.Println("Mutation check: no order, cancellation, approval, wrap, redeem, or transaction was sent")
	return nil
}

func printExchangeReadiness(snapshot *order.ExchangeReadiness) {
	fmt.Printf("Polygon POL: %s\n", formatUnits(snapshot.POLBalance, 18, 6))
	fmt.Printf("Polygon USDC.e: %s\n", formatUnits(snapshot.USDCeBalance, 6, 6))
	fmt.Printf("Polygon pUSD: %s\n", formatUnits(snapshot.PUSDBalance, 6, 6))
	fmt.Printf("USDC.e -> onramp allowance: %v\n", snapshot.USDCeOnrampAllowance.Sign() > 0)
	fmt.Printf("pUSD -> CTF exchange allowance: %v\n", snapshot.PUSDCTFExchangeAllowance.Sign() > 0)
	fmt.Printf("pUSD -> NegRisk exchange allowance: %v\n", snapshot.PUSDNegRiskExchangeAllowance.Sign() > 0)
	fmt.Printf("pUSD -> NegRisk adapter allowance: %v\n", snapshot.PUSDNegRiskAdapterAllowance.Sign() > 0)
	fmt.Printf("CTF sell approval: %v\n", snapshot.CTFExchangeApproved)
	fmt.Printf("NegRisk sell approval: %v\n", snapshot.NegRiskExchangeApproved)
	fmt.Printf("Standard BUY path funded and approved: %v\n", snapshot.PUSDBalance.Sign() > 0 && snapshot.PUSDCTFExchangeAllowance.Sign() > 0)
}

func formatUnits(raw *big.Int, decimals, precision int) string {
	if raw == nil {
		return "0"
	}
	value := new(big.Float).SetInt(raw)
	value.Quo(value, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))
	return value.Text('f', precision)
}

func saveBuyTime(asset string) {
	path := "db/buy_times.json"
	data := map[string]string{}
	if raw, err := os.ReadFile(path); err == nil {
		json.Unmarshal(raw, &data)
	}
	data[asset] = time.Now().Format(time.RFC3339)
	if out, err := json.MarshalIndent(data, "", "  "); err == nil {
		os.WriteFile(path, out, 0644)
		slog.Info("buy_time_saved", "asset", asset[:20])
	}
}

type dataAPIPosition struct {
	Size         float64 `json:"size"`
	CurPrice     float64 `json:"curPrice"`
	CurrentValue float64 `json:"currentValue"`
	Title        string  `json:"title"`
	Outcome      string  `json:"outcome"`
	Asset        string  `json:"asset"`
	ConditionID  string  `json:"conditionId"`
	Redeemable   bool    `json:"redeemable"`
	NegativeRisk bool    `json:"negativeRisk"`
	OutcomeIndex int     `json:"outcomeIndex"`
}

func runRedeemAll(oc *order.OnChain, walletAddr string) {
	reqURL := "https://data-api.polymarket.com/positions?user=" + strings.ToLower(walletAddr) + "&sizeThreshold=0.01&limit=200"
	resp, err := nethttp.Get(reqURL)
	if err != nil {
		slog.Error("data_api_failed", "err", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 400 {
		slog.Error("data_api_error", "status", resp.StatusCode, "body", string(body[:min(len(body), 200)]))
		os.Exit(1)
	}
	var positions []dataAPIPosition
	if err := json.Unmarshal(body, &positions); err != nil {
		slog.Error("data_api_decode", "err", err)
		os.Exit(1)
	}

	redeemed := map[string]bool{}
	if raw, err := os.ReadFile("db/redeemed.json"); err == nil {
		json.Unmarshal(raw, &redeemed)
	}

	var toRedeem []dataAPIPosition
	for _, p := range positions {
		if !p.Redeemable || p.Size < 0.01 {
			continue
		}
		if redeemValue(p) <= 0 {
			slog.Info("redeem_skip_zero_value", "asset", p.Asset[:min(len(p.Asset), 20)], "title", p.Title[:min(len(p.Title), 40)])
			continue
		}
		if redeemed[p.Asset] {
			slog.Info("already_redeemed", "asset", p.Asset[:20], "title", p.Title[:min(len(p.Title), 40)])
			continue
		}
		toRedeem = append(toRedeem, p)
	}

	if len(toRedeem) == 0 {
		fmt.Println("No positions to redeem.")
		bal, err := oc.PUSDBalance(context.Background())
		if err == nil {
			f, _ := new(big.Float).Quo(new(big.Float).SetInt(bal), new(big.Float).SetFloat64(1e6)).Float64()
			fmt.Printf("pUSD balance: $%.2f\n", f)
		}
		return
	}

	fmt.Printf("Found %d positions to redeem:\n", len(toRedeem))
	for _, p := range toRedeem {
		val := redeemValue(p)
		fmt.Printf("  %s · %s · %.1f shares · cur=$%.3f · val=$%.2f · neg=%v\n",
			p.Title[:min(len(p.Title), 50)], p.Outcome, p.Size, p.CurPrice, val, p.NegativeRisk)
	}
	fmt.Println()

	redeemed_count := 0
	for _, p := range toRedeem {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		err := oc.RedeemPosition(ctx, p.ConditionID, p.OutcomeIndex, p.Size, p.NegativeRisk)
		cancel()
		if err != nil {
			slog.Error("redeem_failed", "title", p.Title[:min(len(p.Title), 40)], "err", err)
			continue
		}
		redeemed[p.Asset] = true
		redeemed_count++
		val := redeemValue(p)
		fmt.Printf("✅ Redeemed: %s · %s · $%.2f\n", p.Title[:min(len(p.Title), 50)], p.Outcome, val)
	}

	if out, err := json.MarshalIndent(redeemed, "", "  "); err == nil {
		os.WriteFile("db/redeemed.json", out, 0644)
	}

	bal, err := oc.PUSDBalance(context.Background())
	if err == nil {
		f, _ := new(big.Float).Quo(new(big.Float).SetInt(bal), new(big.Float).SetFloat64(1e6)).Float64()
		fmt.Printf("\nRedeemed %d positions. pUSD balance: $%.2f\n", redeemed_count, f)
	}
}

func redeemValue(p dataAPIPosition) float64 {
	if p.CurrentValue > 0 {
		return p.CurrentValue
	}
	return p.Size * p.CurPrice
}
