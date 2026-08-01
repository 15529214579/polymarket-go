package order

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	PolygonRPC            = "https://polygon-bor-rpc.publicnode.com"
	USDCeAddress          = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"
	PUSDAddress           = "0xC011a7E12a19f7B1f670d46F03B03f3342E82DFB"
	CollateralOnrampAddr  = "0x93070a847efEf7F70739046A929D47a521F5B8ee"
	ConditionalTokensAddr = "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
	NegRiskAdapterAddr    = "0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"
	USDCDecimals          = 6
	ApproveGas            = 60_000
	WrapGas               = 150_000
	RedeemGas             = 300_000
)

var (
	erc20ABI         abi.ABI
	erc1155ABI       abi.ABI
	onrampABI        abi.ABI
	ctfRedeemABI     abi.ABI
	negRiskRedeemABI abi.ABI
)

func init() {
	var err error
	erc20ABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
		{"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}
	]`))
	if err != nil {
		panic("order: parse erc20 abi: " + err.Error())
	}
	erc1155ABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[{"name":"account","type":"address"},{"name":"id","type":"uint256"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"name":"account","type":"address"},{"name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"name":"","type":"bool"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"name":"operator","type":"address"},{"name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"stateMutability":"nonpayable","type":"function"}
	]`))
	if err != nil {
		panic("order: parse erc1155 abi: " + err.Error())
	}

	onrampABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[{"name":"_asset","type":"address"},{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"wrap","outputs":[],"stateMutability":"nonpayable","type":"function"}
	]`))
	if err != nil {
		panic("order: parse onramp abi: " + err.Error())
	}

	ctfRedeemABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[{"name":"collateralToken","type":"address"},{"name":"parentCollectionId","type":"bytes32"},{"name":"conditionId","type":"bytes32"},{"name":"indexSets","type":"uint256[]"}],"name":"redeemPositions","outputs":[],"stateMutability":"nonpayable","type":"function"},
		{"inputs":[{"name":"parentCollectionId","type":"bytes32"},{"name":"conditionId","type":"bytes32"},{"name":"indexSet","type":"uint256"}],"name":"getCollectionId","outputs":[{"name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"},
		{"inputs":[{"name":"collateralToken","type":"address"},{"name":"collectionId","type":"bytes32"}],"name":"getPositionId","outputs":[{"name":"","type":"uint256"}],"stateMutability":"pure","type":"function"}
	]`))
	if err != nil {
		panic("order: parse ctf redeem abi: " + err.Error())
	}

	negRiskRedeemABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[{"name":"conditionId","type":"bytes32"},{"name":"amounts","type":"uint256[]"}],"name":"redeemPositions","outputs":[],"stateMutability":"nonpayable","type":"function"}
	]`))
	if err != nil {
		panic("order: parse neg-risk redeem abi: " + err.Error())
	}
}

type OnChain struct {
	client  *ethclient.Client
	privKey *ecdsa.PrivateKey
	address common.Address
	chainID *big.Int
}

type ExchangeReadiness struct {
	POLBalance                   *big.Int
	USDCeBalance                 *big.Int
	PUSDBalance                  *big.Int
	USDCeOnrampAllowance         *big.Int
	PUSDCTFExchangeAllowance     *big.Int
	PUSDNegRiskExchangeAllowance *big.Int
	PUSDNegRiskAdapterAllowance  *big.Int
	CTFExchangeApproved          bool
	NegRiskExchangeApproved      bool
}

func NewOnChain(rpcURL string, wallet *Wallet) (*OnChain, error) {
	client, err := dialPolygon(rpcURL)
	if err != nil {
		return nil, err
	}
	return &OnChain{
		client:  client,
		privKey: wallet.privKey,
		address: wallet.address,
		chainID: big.NewInt(PolygonChainID),
	}, nil
}

// NewReadOnlyOnChain creates a client that can only perform eth_call and
// balance queries. Transaction submission fails because no private key exists.
func NewReadOnlyOnChain(rpcURL string, address common.Address) (*OnChain, error) {
	client, err := dialPolygon(rpcURL)
	if err != nil {
		return nil, err
	}
	return &OnChain{
		client:  client,
		address: address,
		chainID: big.NewInt(PolygonChainID),
	}, nil
}

func dialPolygon(rpcURL string) (*ethclient.Client, error) {
	if rpcURL == "" {
		rpcURL = PolygonRPC
	}
	directHTTP := &http.Client{
		Transport: &http.Transport{},
		Timeout:   30 * time.Second,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rpcClient, err := rpc.DialOptions(ctx, rpcURL, rpc.WithHTTPClient(directHTTP))
	if err != nil {
		return nil, fmt.Errorf("order: dial polygon rpc: %w", err)
	}
	return ethclient.NewClient(rpcClient), nil
}

func (o *OnChain) ExchangeReadiness(ctx context.Context) (*ExchangeReadiness, error) {
	pol, err := o.client.BalanceAt(ctx, o.address, nil)
	if err != nil {
		return nil, fmt.Errorf("POL balance: %w", err)
	}
	usdce, err := o.USDCeBalance(ctx)
	if err != nil {
		return nil, err
	}
	pusd, err := o.PUSDBalance(ctx)
	if err != nil {
		return nil, err
	}
	usdceOnramp, err := o.getAllowance(ctx, common.HexToAddress(USDCeAddress), common.HexToAddress(CollateralOnrampAddr))
	if err != nil {
		return nil, fmt.Errorf("USDC.e onramp allowance: %w", err)
	}
	pusdToken := common.HexToAddress(PUSDAddress)
	pusdCTF, err := o.getAllowance(ctx, pusdToken, common.HexToAddress(V2ExchangeAddress))
	if err != nil {
		return nil, fmt.Errorf("pUSD CTF exchange allowance: %w", err)
	}
	pusdNegRisk, err := o.getAllowance(ctx, pusdToken, common.HexToAddress(V2NegRiskExchangeAddress))
	if err != nil {
		return nil, fmt.Errorf("pUSD NegRisk exchange allowance: %w", err)
	}
	pusdAdapter, err := o.getAllowance(ctx, pusdToken, common.HexToAddress(NegRiskAdapterAddr))
	if err != nil {
		return nil, fmt.Errorf("pUSD NegRisk adapter allowance: %w", err)
	}
	ctf := common.HexToAddress(ConditionalTokensAddr)
	ctfApproved, err := o.isERC1155Approved(ctx, ctf, common.HexToAddress(V2ExchangeAddress))
	if err != nil {
		return nil, err
	}
	negRiskApproved, err := o.isERC1155Approved(ctx, ctf, common.HexToAddress(V2NegRiskExchangeAddress))
	if err != nil {
		return nil, err
	}
	return &ExchangeReadiness{
		POLBalance:                   pol,
		USDCeBalance:                 usdce,
		PUSDBalance:                  pusd,
		USDCeOnrampAllowance:         usdceOnramp,
		PUSDCTFExchangeAllowance:     pusdCTF,
		PUSDNegRiskExchangeAllowance: pusdNegRisk,
		PUSDNegRiskAdapterAllowance:  pusdAdapter,
		CTFExchangeApproved:          ctfApproved,
		NegRiskExchangeApproved:      negRiskApproved,
	}, nil
}

func (o *OnChain) USDCeBalance(ctx context.Context) (*big.Int, error) {
	return o.erc20Balance(ctx, common.HexToAddress(USDCeAddress), "USDC.e")
}

func (o *OnChain) PUSDBalance(ctx context.Context) (*big.Int, error) {
	return o.erc20Balance(ctx, common.HexToAddress(PUSDAddress), "pUSD")
}

func (o *OnChain) erc20Balance(ctx context.Context, token common.Address, name string) (*big.Int, error) {
	data, err := erc20ABI.Pack("balanceOf", o.address)
	if err != nil {
		return nil, err
	}
	result, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("balanceOf %s: %w", name, err)
	}
	out, err := erc20ABI.Unpack("balanceOf", result)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

func (o *OnChain) ConditionalTokenBalance(ctx context.Context, assetID string) (*big.Int, error) {
	tokenID, ok := new(big.Int).SetString(assetID, 10)
	if !ok || tokenID.Sign() < 0 {
		return nil, fmt.Errorf("invalid conditional token ID %q", assetID)
	}
	data, err := erc1155ABI.Pack("balanceOf", o.address, tokenID)
	if err != nil {
		return nil, err
	}
	ctf := common.HexToAddress(ConditionalTokensAddr)
	result, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &ctf, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("conditional token balanceOf: %w", err)
	}
	out, err := erc1155ABI.Unpack("balanceOf", result)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

func (o *OnChain) WrapAll(ctx context.Context) error {
	bal, err := o.USDCeBalance(ctx)
	if err != nil {
		return fmt.Errorf("check USDC.e balance: %w", err)
	}
	if bal.Sign() <= 0 {
		slog.Info("wrap_skip", "reason", "zero USDC.e balance")
		return nil
	}
	slog.Info("wrap_start", "usdc_e_raw", bal.String(), "usdc_e", formatUSDC(bal))

	onramp := common.HexToAddress(CollateralOnrampAddr)
	allowance, err := o.getAllowance(ctx, common.HexToAddress(USDCeAddress), onramp)
	if err != nil {
		return fmt.Errorf("check allowance: %w", err)
	}

	if allowance.Cmp(bal) < 0 {
		maxApproval := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
		slog.Info("wrap_approve", "current_allowance", allowance.String(), "needed", bal.String())
		if err := o.approve(ctx, common.HexToAddress(USDCeAddress), onramp, maxApproval); err != nil {
			return fmt.Errorf("approve USDC.e: %w", err)
		}
		slog.Info("wrap_approved")
	}

	slog.Info("wrap_calling", "amount", formatUSDC(bal))
	if err := o.wrap(ctx, bal); err != nil {
		return fmt.Errorf("wrap: %w", err)
	}

	pusd, err := o.PUSDBalance(ctx)
	if err != nil {
		slog.Warn("wrap_verify_failed", "err", err)
	} else {
		slog.Info("wrap_done", "pusd_balance", formatUSDC(pusd))
	}
	return nil
}

func (o *OnChain) ApproveExchanges(ctx context.Context) error {
	pusd := common.HexToAddress(PUSDAddress)
	maxApproval := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	exchanges := []struct {
		name string
		addr string
	}{
		{"CTFExchange", V2ExchangeAddress},
		{"NegRiskExchange", V2NegRiskExchangeAddress},
		{"NegRiskAdapter", "0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"},
	}

	for _, ex := range exchanges {
		spender := common.HexToAddress(ex.addr)
		allowance, err := o.getAllowance(ctx, pusd, spender)
		if err != nil {
			return fmt.Errorf("check pUSD allowance for %s: %w", ex.name, err)
		}
		if allowance.Sign() > 0 {
			slog.Info("exchange_allowance_ok", "exchange", ex.name, "allowance", allowance.String())
			continue
		}
		slog.Info("exchange_approve", "exchange", ex.name, "spender", ex.addr)
		if err := o.approve(ctx, pusd, spender, maxApproval); err != nil {
			return fmt.Errorf("approve pUSD for %s: %w", ex.name, err)
		}
		slog.Info("exchange_approved", "exchange", ex.name)
	}

	ctf := common.HexToAddress(ConditionalTokensAddr)
	erc1155Operators := []struct {
		name string
		addr common.Address
	}{
		{"CTFExchange", common.HexToAddress(V2ExchangeAddress)},
		{"NegRiskExchange", common.HexToAddress(V2NegRiskExchangeAddress)},
	}
	for _, op := range erc1155Operators {
		if err := o.ensureERC1155Approval(ctx, ctf, op.addr); err != nil {
			return fmt.Errorf("ERC1155 approve for %s: %w", op.name, err)
		}
	}

	return nil
}

func (o *OnChain) getAllowance(ctx context.Context, token, spender common.Address) (*big.Int, error) {
	data, err := erc20ABI.Pack("allowance", o.address, spender)
	if err != nil {
		return nil, err
	}
	result, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	out, err := erc20ABI.Unpack("allowance", result)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

func (o *OnChain) approve(ctx context.Context, token, spender common.Address, amount *big.Int) error {
	data, err := erc20ABI.Pack("approve", spender, amount)
	if err != nil {
		return err
	}
	return o.sendTx(ctx, token, data, ApproveGas)
}

func (o *OnChain) wrap(ctx context.Context, amount *big.Int) error {
	asset := common.HexToAddress(USDCeAddress)
	data, err := onrampABI.Pack("wrap", asset, o.address, amount)
	if err != nil {
		return err
	}
	return o.sendTx(ctx, common.HexToAddress(CollateralOnrampAddr), data, WrapGas)
}

func (o *OnChain) sendTx(ctx context.Context, to common.Address, data []byte, gasLimit uint64) error {
	if o.privKey == nil {
		return fmt.Errorf("transaction signing is unavailable on a read-only client")
	}
	nonce, err := o.client.PendingNonceAt(ctx, o.address)
	if err != nil {
		return fmt.Errorf("get nonce: %w", err)
	}

	gasTip, err := o.client.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("gas tip: %w", err)
	}
	head, err := o.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("get header: %w", err)
	}
	baseFee := head.BaseFee
	maxFee := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTip)

	estGas, err := o.client.EstimateGas(ctx, ethereum.CallMsg{
		From:  o.address,
		To:    &to,
		Data:  data,
		Value: big.NewInt(0),
	})
	if err != nil {
		slog.Warn("gas_estimate_failed", "err", err, "using_default", gasLimit)
	} else {
		gasLimit = estGas * 120 / 100
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   o.chainID,
		Nonce:     nonce,
		GasTipCap: gasTip,
		GasFeeCap: maxFee,
		Gas:       gasLimit,
		To:        &to,
		Value:     big.NewInt(0),
		Data:      data,
	})
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(o.chainID), o.privKey)
	if err != nil {
		return fmt.Errorf("sign tx: %w", err)
	}

	if err := o.client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("send tx: %w", err)
	}

	slog.Info("tx_sent", "hash", signedTx.Hash().Hex(), "to", to.Hex(), "gas", gasLimit)

	receipt, err := waitForReceipt(ctx, o.client, signedTx.Hash(), 120*time.Second)
	if err != nil {
		return fmt.Errorf("wait receipt: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("tx reverted: %s", signedTx.Hash().Hex())
	}
	slog.Info("tx_confirmed", "hash", signedTx.Hash().Hex(), "gas_used", receipt.GasUsed)
	return nil
}

func waitForReceipt(ctx context.Context, client *ethclient.Client, hash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	deadline := time.After(timeout)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-deadline:
			return nil, fmt.Errorf("timeout waiting for tx %s", hash.Hex())
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, hash)
			if err != nil {
				continue
			}
			return receipt, nil
		}
	}
}

func formatUSDC(raw *big.Int) string {
	f := new(big.Float).SetInt(raw)
	divisor := new(big.Float).SetFloat64(1e6)
	f.Quo(f, divisor)
	return f.Text('f', 6)
}

func (o *OnChain) RedeemPosition(ctx context.Context, conditionIDHex, assetID string, outcomeIndex int, size float64, negRisk bool) error {
	conditionBytes := common.FromHex(conditionIDHex)
	if len(conditionBytes) != common.HashLength {
		return fmt.Errorf("invalid condition ID %q", conditionIDHex)
	}
	conditionID := common.BytesToHash(conditionBytes)
	indexSet := new(big.Int).Lsh(big.NewInt(1), uint(outcomeIndex))
	collateral, err := o.positionCollateral(ctx, conditionID, indexSet, assetID)
	if err != nil {
		return err
	}
	tokenBefore, err := o.ConditionalTokenBalance(ctx, assetID)
	if err != nil {
		return err
	}
	if tokenBefore.Sign() <= 0 {
		return fmt.Errorf("conditional token %s has zero on-chain balance", assetID)
	}
	collateralName := "pUSD"
	if collateral == common.HexToAddress(USDCeAddress) {
		collateralName = "USDC.e"
	}
	collateralBefore, err := o.erc20Balance(ctx, collateral, collateralName)
	if err != nil {
		return err
	}

	slog.Info("redeem_start",
		"conditionId", conditionIDHex,
		"outcomeIndex", outcomeIndex,
		"size", size,
		"negRisk", negRisk,
		"collateral", collateralName,
	)

	if negRisk {
		if collateral != common.HexToAddress(PUSDAddress) {
			return fmt.Errorf("legacy USDC.e NegRisk redemption is not supported by the configured V2 adapter")
		}
		// NegRisk: first ensure CTF ERC1155 approval for NegRiskAdapter, then try adapter.
		// If adapter fails, fall back to CTF direct redeem.
		if err := o.ensureERC1155Approval(ctx, common.HexToAddress(ConditionalTokensAddr), common.HexToAddress(NegRiskAdapterAddr)); err != nil {
			slog.Warn("erc1155_approval_failed", "err", err)
		}
		sizeRaw := new(big.Int).SetInt64(int64(size * 1e6))
		amounts := []*big.Int{big.NewInt(0), big.NewInt(0)}
		if outcomeIndex < len(amounts) {
			amounts[outcomeIndex] = sizeRaw
		}
		data, err := negRiskRedeemABI.Pack("redeemPositions", conditionID, amounts)
		if err != nil {
			return fmt.Errorf("pack neg-risk redeem: %w", err)
		}
		if err := o.sendTx(ctx, common.HexToAddress(NegRiskAdapterAddr), data, RedeemGas); err == nil {
			return o.verifyRedemption(ctx, assetID, collateral, collateralName, tokenBefore, collateralBefore)
		} else {
			slog.Warn("negrisk_adapter_redeem_failed_fallback_ctf", "err", err)
		}
	}

	// CTF direct redeem (works for non-NegRisk; fallback for NegRisk)
	data, err := ctfRedeemABI.Pack("redeemPositions",
		collateral,
		common.Hash{},
		conditionID,
		[]*big.Int{indexSet},
	)
	if err != nil {
		return fmt.Errorf("pack ctf redeem: %w", err)
	}
	if err := o.sendTx(ctx, common.HexToAddress(ConditionalTokensAddr), data, RedeemGas); err != nil {
		return err
	}
	return o.verifyRedemption(ctx, assetID, collateral, collateralName, tokenBefore, collateralBefore)
}

func (o *OnChain) positionCollateral(ctx context.Context, conditionID common.Hash, indexSet *big.Int, assetID string) (common.Address, error) {
	target, ok := new(big.Int).SetString(assetID, 10)
	if !ok || target.Sign() < 0 {
		return common.Address{}, fmt.Errorf("invalid conditional token ID %q", assetID)
	}
	ctf := common.HexToAddress(ConditionalTokensAddr)
	collectionData, err := ctfRedeemABI.Pack("getCollectionId", common.Hash{}, conditionID, indexSet)
	if err != nil {
		return common.Address{}, err
	}
	collectionRaw, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &ctf, Data: collectionData}, nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("getCollectionId: %w", err)
	}
	collectionOut, err := ctfRedeemABI.Unpack("getCollectionId", collectionRaw)
	if err != nil {
		return common.Address{}, err
	}
	collectionID := collectionOut[0].([32]byte)
	for _, candidate := range []common.Address{common.HexToAddress(PUSDAddress), common.HexToAddress(USDCeAddress)} {
		positionData, err := ctfRedeemABI.Pack("getPositionId", candidate, collectionID)
		if err != nil {
			return common.Address{}, err
		}
		positionRaw, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &ctf, Data: positionData}, nil)
		if err != nil {
			return common.Address{}, fmt.Errorf("getPositionId: %w", err)
		}
		positionOut, err := ctfRedeemABI.Unpack("getPositionId", positionRaw)
		if err != nil {
			return common.Address{}, err
		}
		if positionOut[0].(*big.Int).Cmp(target) == 0 {
			return candidate, nil
		}
	}
	return common.Address{}, fmt.Errorf("conditional token %s does not match pUSD or USDC.e collateral", assetID)
}

func (o *OnChain) verifyRedemption(ctx context.Context, assetID string, collateral common.Address, collateralName string, tokenBefore, collateralBefore *big.Int) error {
	tokenAfter, err := o.ConditionalTokenBalance(ctx, assetID)
	if err != nil {
		return fmt.Errorf("verify conditional token balance: %w", err)
	}
	collateralAfter, err := o.erc20Balance(ctx, collateral, collateralName)
	if err != nil {
		return fmt.Errorf("verify collateral balance: %w", err)
	}
	if tokenAfter.Cmp(tokenBefore) >= 0 {
		return fmt.Errorf("redemption confirmed but conditional token balance did not decrease")
	}
	if collateralAfter.Cmp(collateralBefore) <= 0 {
		return fmt.Errorf("redemption confirmed but %s balance did not increase", collateralName)
	}
	slog.Info("redeem_verified",
		"asset", assetID,
		"token_burned", new(big.Int).Sub(tokenBefore, tokenAfter).String(),
		"collateral", collateralName,
		"collateral_received", new(big.Int).Sub(collateralAfter, collateralBefore).String(),
	)
	return nil
}

func (o *OnChain) ensureERC1155Approval(ctx context.Context, tokenContract, operator common.Address) error {
	approved, err := o.isERC1155Approved(ctx, tokenContract, operator)
	if err != nil {
		return err
	}
	if approved {
		slog.Info("erc1155_already_approved", "token", tokenContract.Hex(), "operator", operator.Hex())
		return nil
	}

	slog.Info("erc1155_approving", "token", tokenContract.Hex(), "operator", operator.Hex())
	setData, err := erc1155ABI.Pack("setApprovalForAll", operator, true)
	if err != nil {
		return err
	}
	return o.sendTx(ctx, tokenContract, setData, ApproveGas)
}

func (o *OnChain) isERC1155Approved(ctx context.Context, tokenContract, operator common.Address) (bool, error) {
	data, err := erc1155ABI.Pack("isApprovedForAll", o.address, operator)
	if err != nil {
		return false, err
	}
	result, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &tokenContract, Data: data}, nil)
	if err != nil {
		return false, fmt.Errorf("isApprovedForAll %s: %w", operator.Hex(), err)
	}
	out, err := erc1155ABI.Unpack("isApprovedForAll", result)
	if err != nil {
		return false, err
	}
	return out[0].(bool), nil
}
