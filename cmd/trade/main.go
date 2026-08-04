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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/order"
	"github.com/ethereum/go-ethereum/common"
)

const projectWalletAddress = "0x015282e9b720E072A9B87eEeaE738C6Bb039Bd9e"

func main() {
	assetID := flag.String("asset", "", "ERC1155 token ID (YES side)")
	market := flag.String("market", "", "conditionID for logging")
	sizeUSD := flag.Float64("size", 0, "order size in USDC")
	limitPx := flag.Float64("price", 0, "limit price (0..1)")
	side := flag.String("side", "BUY", "BUY or SELL")
	negRisk := flag.Bool("negrisk", false, "use NegRisk exchange address")
	dryRun := flag.Bool("dry-run", false, "validate and print intent locally; no network or wallet mutation")
	wrapApprove := flag.Bool("wrap-approve", false, "wallet maintenance: wrap all USDC.e, approve exchanges, then exit")
	rpcURL := flag.String("rpc", "", "Polygon RPC URL (default: polygon-bor-rpc.publicnode.com)")
	queryOpenOrders := flag.Bool("open-orders", false, "query CLOB open orders without cancelling")
	cancelOpenOrders := flag.Bool("cancel-open-orders", false, "explicitly cancel every open CLOB order")
	readiness := flag.Bool("readiness", false, "read-only wallet, L1/L2 auth, balance, allowance, and open-order checks")
	publicReadiness := flag.Bool("readiness-public", false, "read-only public CLOB and Polygon checks; no private key required")
	expectedAddress := flag.String("expected-address", projectWalletAddress, "wallet address that must match before authenticated operations")
	redeemAll := flag.Bool("redeem-all", false, "redeem all redeemable positions and exit")
	redeemedStatePath := flag.String("redeemed-state", "db/live/redeemed.json", "live redemption state path; only used with -redeem-all")
	liveArmFile := flag.String("live-arm-file", "db/live-trading.enabled", "short-lived, wallet-bound live trading arm file")
	liveDisableFile := flag.String("live-disable-file", "db/live-trading.disabled", "live trading kill-switch file")
	liveMaxOrderUSD := flag.Float64("live-max-order-usd", 20.0, "hard maximum notional for one live BUY; exits are not capped")
	liveMaxSessionBuyUSD := flag.Float64("live-max-session-buy-usd", 100.0, "hard maximum reserved/filled BUY notional for one arm window")
	liveMaxArmDuration := flag.Duration("live-max-arm-duration", 24*time.Hour, "maximum accepted live arm-file validity window")
	liveSessionStatePath := flag.String("live-session-state", "db/live/live-session.json", "durable live guard session state")
	executionLedgerPath := flag.String("execution-ledger", "db/live/orders.sqlite", "durable live order execution ledger")
	buyTimesStatePath := flag.String("buy-times-state", "db/live/buy_times.json", "durable live buy-time state")
	flag.Parse()

	if err := validateOperationFlags(*wrapApprove, *redeemAll, *readiness, *publicReadiness, *queryOpenOrders, *cancelOpenOrders, *dryRun); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
	managementMode := *queryOpenOrders || *cancelOpenOrders || *readiness || *publicReadiness || *redeemAll || *wrapApprove
	if !managementMode {
		for _, state := range []struct {
			name string
			path *string
		}{
			{"live session", liveSessionStatePath},
			{"execution ledger", executionLedgerPath},
			{"buy times", buyTimesStatePath},
		} {
			validated, err := validateLiveRuntimeStatePath(*state.path, state.name)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1)
			}
			*state.path = validated
		}
	}
	if !managementMode && (*assetID == "" || *limitPx <= 0) {
		fmt.Fprintln(os.Stderr, "Usage: trade -asset <tokenID> -price <0..1> -size <usd> [-negrisk] [-dry-run]")
		fmt.Fprintln(os.Stderr, "Maintenance: trade -wrap-approve | -redeem-all | -readiness | -readiness-public")
		os.Exit(1)
	}
	if *expectedAddress != "" && !common.IsHexAddress(*expectedAddress) {
		fmt.Fprintf(os.Stderr, "ERROR: invalid -expected-address %q\n", *expectedAddress)
		os.Exit(1)
	}

	slog.Info("trade_init", "asset", *assetID, "size", *sizeUSD, "price", *limitPx, "side", *side, "negrisk", *negRisk)

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

	wallet, err := order.LoadWalletFromBitwarden("Polymarket-Go Wallet", "mnemonic", "")
	if err != nil {
		slog.Error("wallet_load_failed", "err", err)
		os.Exit(1)
	}
	slog.Info("wallet_loaded", "address", wallet.Address().Hex())
	if *expectedAddress != "" && wallet.Address() != common.HexToAddress(*expectedAddress) {
		slog.Error("wallet_address_mismatch", "loaded", wallet.Address().Hex(), "expected", common.HexToAddress(*expectedAddress).Hex())
		os.Exit(1)
	}
	if *wrapApprove {
		oc, err := order.NewOnChain(*rpcURL, wallet)
		if err != nil {
			slog.Error("onchain_init_failed", "err", err)
			os.Exit(1)
		}
		if err := runWrapApprove(oc); err != nil {
			slog.Error("wrap_approve_failed", "err", err)
			os.Exit(1)
		}
		return
	}
	if *redeemAll {
		oc, err := order.NewOnChain(*rpcURL, wallet)
		if err != nil {
			slog.Error("onchain_init_failed", "err", err)
			os.Exit(1)
		}
		if err := runRedeemAll(oc, wallet.Address().Hex(), *redeemedStatePath); err != nil {
			slog.Error("redeem_all_failed", "err", err)
			os.Exit(1)
		}
		return
	}
	order.InitProxy()
	var liveGuardCfg order.LiveGuardConfig
	if !managementMode {
		liveGuardCfg = order.LiveGuardConfig{
			ArmFile:          *liveArmFile,
			DisableFile:      *liveDisableFile,
			ExpectedWallet:   wallet.Address().Hex(),
			MaxOrderUSD:      *liveMaxOrderUSD,
			MaxSessionBuyUSD: *liveMaxSessionBuyUSD,
			MaxArmDuration:   *liveMaxArmDuration,
			SessionStatePath: *liveSessionStatePath,
		}
		if err := order.CheckLiveGuard(liveGuardCfg); err != nil {
			slog.Error("live_guard_rejected", "err", err)
			os.Exit(1)
		}
	}

	creds, err := order.DeriveExistingAPIKey(order.ClobBaseURL, wallet)
	if err != nil {
		slog.Error("api_key_derive_failed", "err", err)
		os.Exit(1)
	}
	slog.Info("api_key_derived")
	client := order.NewV2Client(wallet, creds, *negRisk)
	var submitClient order.Client = client
	var executionLedger *order.ExecutionLedger
	if !managementMode {
		executionLedger, err = order.OpenExecutionLedger(*executionLedgerPath)
		if err != nil {
			slog.Error("execution_ledger_init_failed", "err", err)
			os.Exit(1)
		}
		defer func() {
			if err := executionLedger.Close(); err != nil {
				slog.Error("execution_ledger_close_failed", "err", err)
			}
		}()
		if err := executionLedger.ReconcileLive(context.Background(), client); err != nil {
			slog.Warn("execution_ledger_reconcile_partial", "err", err)
		}
		unresolved, err := executionLedger.UnresolvedCount("live")
		if err != nil {
			slog.Error("execution_ledger_count_failed", "err", err)
			os.Exit(1)
		}
		if unresolved > 0 {
			slog.Error("execution_ledger_unresolved", "count", unresolved)
			os.Exit(1)
		}
		if err := applyRecoveredTradeFills(executionLedger, *buyTimesStatePath); err != nil {
			slog.Error("execution_ledger_recovery_failed", "err", err)
			os.Exit(1)
		}
		guardedClient, err := order.NewGuardedClient(client, liveGuardCfg)
		if err != nil {
			slog.Error("live_guard_init_failed", "err", err)
			os.Exit(1)
		}
		if err := guardedClient.CheckReady(); err != nil {
			slog.Error("live_guard_rejected", "err", err)
			os.Exit(1)
		}
		ledgerClient, err := order.NewLedgerClient(guardedClient, executionLedger, "live")
		if err != nil {
			slog.Error("execution_ledger_wrap_failed", "err", err)
			os.Exit(1)
		}
		submitClient = ledgerClient
	}

	if *readiness {
		if err := runAuthenticatedReadiness(*rpcURL, wallet, client); err != nil {
			slog.Error("authenticated_readiness_failed", "err", err)
			os.Exit(1)
		}
		return
	}

	if *queryOpenOrders {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		openOrders, err := client.GetOpenOrders(ctx)
		cancel()
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
		err := client.CancelAllOpen(ctx)
		cancel()
		if err != nil {
			slog.Error("cancel_all_failed", "err", err)
			os.Exit(1)
		}
		fmt.Println("All open orders cancelled.")
		return
	}

	if *sizeUSD <= 0 {
		fmt.Fprintln(os.Stderr, "ERROR: -size is required")
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
	result, err := submitClient.Submit(ctx, intent)
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
		if intent.Side == order.Buy {
			if err := saveBuyTime(*buyTimesStatePath, *assetID, result.FilledAt); err != nil {
				slog.Error("buy_time_save_failed", "asset", *assetID, "err", err)
				os.Exit(1)
			}
		}
		if err := executionLedger.MarkApplied(result.ExecutionID); err != nil {
			slog.Error("execution_ledger_mark_applied_failed", "execution_id", result.ExecutionID, "err", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("\n⚠️  Status: %s — %s\n", result.Status, result.Error)
	}
}

func validateOperationFlags(wrapApprove, redeemAll, readiness, publicReadiness, queryOpenOrders, cancelOpenOrders, dryRun bool) error {
	clobManagement := queryOpenOrders || cancelOpenOrders
	if wrapApprove && (redeemAll || readiness || publicReadiness || clobManagement || dryRun) {
		return fmt.Errorf("-wrap-approve must run as a standalone maintenance action")
	}
	if redeemAll && (readiness || publicReadiness || clobManagement || dryRun) {
		return fmt.Errorf("-redeem-all must run as a standalone maintenance action")
	}
	if readiness && (publicReadiness || clobManagement || dryRun) {
		return fmt.Errorf("-readiness cannot be combined with another operation")
	}
	if publicReadiness && (clobManagement || dryRun) {
		return fmt.Errorf("-readiness-public cannot be combined with another operation")
	}
	if dryRun && clobManagement {
		return fmt.Errorf("-dry-run cannot be combined with CLOB management")
	}
	return nil
}

func runWrapApprove(oc *order.OnChain) error {
	wrapCtx, wrapCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	err := oc.WrapAll(wrapCtx)
	wrapCancel()
	if err != nil {
		return fmt.Errorf("wrap USDC.e: %w", err)
	}

	approveCtx, approveCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	err = oc.ApproveExchanges(approveCtx)
	approveCancel()
	if err != nil {
		return fmt.Errorf("approve exchanges: %w", err)
	}

	pusd, err := oc.PUSDBalance(context.Background())
	if err != nil {
		return fmt.Errorf("check pUSD balance: %w", err)
	}
	pusdFloat, _ := new(big.Float).Quo(
		new(big.Float).SetInt(pusd),
		new(big.Float).SetFloat64(1e6),
	).Float64()
	fmt.Printf("Wallet maintenance complete. pUSD balance: $%.2f\n", pusdFloat)
	return nil
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
		Type:    order.FAK,
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

func applyRecoveredTradeFills(ledger *order.ExecutionLedger, buyTimesPath string) error {
	nonFills, err := ledger.UnappliedNonFills("live")
	if err != nil {
		return err
	}
	for _, record := range nonFills {
		if err := ledger.MarkApplied(record.ID); err != nil {
			return err
		}
		slog.Warn("execution_ledger_nonfill_recovered",
			"execution_id", record.ID,
			"order_id", record.OrderID,
			"status", record.Status,
			"side", record.Intent.Side,
		)
	}
	records, err := ledger.UnappliedFills("live")
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Intent.Side == order.Buy {
			filledAt := record.Result.FilledAt
			if filledAt.IsZero() {
				filledAt = record.UpdatedAt
			}
			if err := saveBuyTime(buyTimesPath, record.Intent.AssetID, filledAt); err != nil {
				return fmt.Errorf("recover execution %s buy time: %w", record.ID, err)
			}
		}
		if err := ledger.MarkApplied(record.ID); err != nil {
			return err
		}
		slog.Warn("execution_ledger_fill_recovered",
			"execution_id", record.ID,
			"order_id", record.Result.OrderID,
			"side", record.Intent.Side,
			"asset", record.Intent.AssetID,
			"filled_size", record.Result.FilledSize,
			"avg_price", record.Result.AvgPrice,
		)
	}
	return nil
}

func saveBuyTime(path, asset string, filledAt time.Time) error {
	if filledAt.IsZero() {
		filledAt = time.Now()
	}
	data := map[string]string{}
	if raw, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	data[asset] = filledAt.UTC().Format(time.RFC3339Nano)
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".buy-times-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(out); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	dirHandle, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	if err := dirHandle.Sync(); err != nil {
		_ = dirHandle.Close()
		return err
	}
	if err := dirHandle.Close(); err != nil {
		return err
	}
	slog.Info("buy_time_saved", "asset", asset[:min(len(asset), 20)])
	return nil
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

func runRedeemAll(oc *order.OnChain, walletAddr, redeemedStatePath string) error {
	var err error
	redeemedStatePath, err = validateLiveRedeemedStatePath(redeemedStatePath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(redeemedStatePath), 0700); err != nil {
		return fmt.Errorf("create redeemed state directory: %w", err)
	}
	httpClient := &nethttp.Client{Timeout: 30 * time.Second}
	fetchCtx, fetchCancel := context.WithTimeout(context.Background(), 45*time.Second)
	positions, err := fetchRedeemablePositions(fetchCtx, httpClient, "https://data-api.polymarket.com", walletAddr)
	fetchCancel()
	if err != nil {
		return err
	}

	redeemed := map[string]bool{}
	if raw, err := os.ReadFile(redeemedStatePath); err == nil {
		_ = json.Unmarshal(raw, &redeemed)
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
		balanceCtx, balanceCancel := context.WithTimeout(context.Background(), 15*time.Second)
		chainBalance, err := oc.ConditionalTokenBalance(balanceCtx, p.Asset)
		balanceCancel()
		if err != nil {
			slog.Warn("redeem_balance_check_failed", "asset", p.Asset[:min(len(p.Asset), 20)], "err", err)
			continue
		}
		if chainBalance.Sign() <= 0 {
			redeemed[p.Asset] = true
			slog.Info("redeem_skip_zero_chain_balance", "asset", p.Asset[:min(len(p.Asset), 20)], "title", p.Title[:min(len(p.Title), 40)])
			continue
		}
		if redeemed[p.Asset] {
			slog.Warn("redeem_local_state_stale", "asset", p.Asset[:min(len(p.Asset), 20)], "onchain_balance", chainBalance.String())
		}
		toRedeem = append(toRedeem, p)
	}

	if len(toRedeem) == 0 {
		fmt.Println("No positions to redeem.")
		balanceCtx, balanceCancel := context.WithTimeout(context.Background(), 15*time.Second)
		bal, err := oc.PUSDBalance(balanceCtx)
		balanceCancel()
		if err == nil {
			f, _ := new(big.Float).Quo(new(big.Float).SetInt(bal), new(big.Float).SetFloat64(1e6)).Float64()
			fmt.Printf("pUSD balance: $%.2f\n", f)
		}
		return nil
	}

	fmt.Printf("Found %d positions to redeem:\n", len(toRedeem))
	for _, p := range toRedeem {
		val := redeemValue(p)
		fmt.Printf("  %s · %s · %.1f shares · cur=$%.3f · val=$%.2f · neg=%v\n",
			p.Title[:min(len(p.Title), 50)], p.Outcome, p.Size, p.CurPrice, val, p.NegativeRisk)
	}
	fmt.Println()

	redeemedCount := 0
	failedCount := 0
	for _, p := range toRedeem {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		err := oc.RedeemPosition(ctx, p.ConditionID, p.Asset, p.OutcomeIndex, p.Size, p.NegativeRisk)
		cancel()
		if err != nil {
			slog.Error("redeem_failed", "title", p.Title[:min(len(p.Title), 40)], "err", err)
			failedCount++
			continue
		}
		redeemed[p.Asset] = true
		redeemedCount++
		val := redeemValue(p)
		fmt.Printf("✅ Redeemed: %s · %s · $%.2f\n", p.Title[:min(len(p.Title), 50)], p.Outcome, val)
	}

	out, err := json.MarshalIndent(redeemed, "", "  ")
	if err != nil {
		return fmt.Errorf("encode redeemed state: %w", err)
	}
	if err := os.WriteFile(redeemedStatePath, out, 0600); err != nil {
		return fmt.Errorf("write redeemed state: %w", err)
	}

	balanceCtx, balanceCancel := context.WithTimeout(context.Background(), 15*time.Second)
	usdce, usdceErr := oc.USDCeBalance(balanceCtx)
	balanceCancel()
	balanceCtx, balanceCancel = context.WithTimeout(context.Background(), 15*time.Second)
	pusd, pusdErr := oc.PUSDBalance(balanceCtx)
	balanceCancel()
	if usdceErr == nil && pusdErr == nil {
		usdceFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(usdce), new(big.Float).SetFloat64(1e6)).Float64()
		pusdFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(pusd), new(big.Float).SetFloat64(1e6)).Float64()
		fmt.Printf("\nRedeemed %d positions. USDC.e balance: $%.2f · pUSD balance: $%.2f\n", redeemedCount, usdceFloat, pusdFloat)
	}
	if failedCount > 0 {
		return fmt.Errorf("%d of %d redemption attempts failed", failedCount, len(toRedeem))
	}
	return nil
}

func fetchRedeemablePositions(ctx context.Context, client *nethttp.Client, baseURL, walletAddr string) ([]dataAPIPosition, error) {
	const pageLimit = 500
	var positions []dataAPIPosition
	for offset := 0; offset <= 10000; offset += pageLimit {
		query := url.Values{
			"user":          {strings.ToLower(strings.TrimSpace(walletAddr))},
			"sizeThreshold": {"0.01"},
			"redeemable":    {"true"},
			"limit":         {strconv.Itoa(pageLimit)},
			"offset":        {strconv.Itoa(offset)},
		}
		req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, strings.TrimRight(baseURL, "/")+"/positions?"+query.Encode(), nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("data API positions: %w", err)
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		closeErr := resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read data API positions: %w", readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("close data API response: %w", closeErr)
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("data API returned HTTP %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
		}
		var page []dataAPIPosition
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("decode data API positions: %w", err)
		}
		positions = append(positions, page...)
		if len(page) < pageLimit {
			return positions, nil
		}
		if offset == 10000 {
			return nil, fmt.Errorf("data API positions exceeded pagination limit")
		}
	}
	return positions, nil
}

func validateLiveRedeemedStatePath(path string) (string, error) {
	return validateLiveRuntimeStatePath(path, "redeemed state")
}

func validateLiveRuntimeStatePath(path, label string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("%s path is required", label)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve %s path: %w", label, err)
	}
	liveRoot, err := filepath.Abs(filepath.Join("db", "live"))
	if err != nil {
		return "", fmt.Errorf("resolve live state root: %w", err)
	}
	rel, err := filepath.Rel(liveRoot, absPath)
	if err != nil {
		return "", fmt.Errorf("compare %s path: %w", label, err)
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s path must be a file under %s", label, liveRoot)
	}
	if filepath.Dir(absPath) != liveRoot {
		return "", fmt.Errorf("%s path must be a direct child of %s", label, liveRoot)
	}
	if err := rejectLiveSymlinkComponents(absPath); err != nil {
		return "", fmt.Errorf("%s path is unsafe: %w", label, err)
	}
	if info, statErr := os.Lstat(liveRoot); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("live state root must not be a symlink")
	} else if statErr != nil && !os.IsNotExist(statErr) {
		return "", fmt.Errorf("inspect live state root: %w", statErr)
	}
	if info, statErr := os.Lstat(absPath); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("%s file must not be a symlink", label)
	} else if statErr != nil && !os.IsNotExist(statErr) {
		return "", fmt.Errorf("inspect %s path: %w", label, statErr)
	}
	return absPath, nil
}

func rejectLiveSymlinkComponents(path string) error {
	for current := path; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		switch {
		case err == nil && info.Mode()&os.ModeSymlink != 0:
			return fmt.Errorf("%s is a symlink", current)
		case err != nil && !os.IsNotExist(err):
			return fmt.Errorf("inspect %s: %w", current, err)
		}
		parent := filepath.Dir(current)
		if parent == current {
			return nil
		}
	}
}

func redeemValue(p dataAPIPosition) float64 {
	if p.CurrentValue > 0 {
		return p.CurrentValue
	}
	return p.Size * p.CurPrice
}
