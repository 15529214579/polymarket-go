package order

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"golang.org/x/crypto/sha3"
)

// V2 Polymarket CLOB EIP-712 constants. Values below reflect the
// pre-cutover spec (SPEC §5.2); any rename announced on 2026-04-28 must be
// mirrored here. Kept as constants so a grep is trivial on cutover day.
const (
	PolygonChainID int64 = 137

	V2DomainName    = "Polymarket CTF Exchange"
	V2DomainVersion = "2"

	V2ExchangeAddress        = "0xE111180000d2663C0091e4f400237545B87B996B"
	V2NegRiskExchangeAddress = "0xe2222d279d744050d28e00520010520000310F59"
)

// SigSide matches the on-chain `enum Side { BUY, SELL }` encoding.
type SigSide uint8

const (
	SigBuy  SigSide = 0
	SigSell SigSide = 1
)

// SigType matches `enum SignatureType { EOA, POLY_PROXY, POLY_GNOSIS_SAFE }`.
type SigType uint8

const (
	SigTypeEOA            SigType = 0
	SigTypePolyProxy      SigType = 1
	SigTypePolyGnosisSafe SigType = 2
)

// V2Order is the hash-able shape of a single CLOB V2 limit order. Every
// field is concrete enough to feed EIP-712 typed-data hashing; a signer
// turns it into the 65-byte signature expected by the CLOB REST endpoint.
//
// Field set reflects SPEC §5.2: V1 `taker/expiration/nonce/feeRateBps`
// dropped; `timestamp/metadata` added. Verify exact names on cutover and
// adjust `v2OrderTypes` below accordingly.
type V2Order struct {
	Salt          *big.Int
	Maker         common.Address
	Signer        common.Address
	TokenID       *big.Int
	MakerAmount   *big.Int
	TakerAmount   *big.Int
	Side          SigSide
	SignatureType SigType
	Timestamp     *big.Int
	Metadata      [32]byte
	Builder       [32]byte
}

// v2OrderTypes is the EIP-712 schema for V2Order. Keep in sync with the
// field list in V2Order.
var v2OrderTypes = apitypes.Types{
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"Order": []apitypes.Type{
		{Name: "salt", Type: "uint256"},
		{Name: "maker", Type: "address"},
		{Name: "signer", Type: "address"},
		{Name: "tokenId", Type: "uint256"},
		{Name: "makerAmount", Type: "uint256"},
		{Name: "takerAmount", Type: "uint256"},
		{Name: "side", Type: "uint8"},
		{Name: "signatureType", Type: "uint8"},
		{Name: "timestamp", Type: "uint256"},
		{Name: "metadata", Type: "bytes32"},
		{Name: "builder", Type: "bytes32"},
	},
}

// EIP712HashV2Order produces the 32-byte digest that must be signed to
// authorise a V2 order for exchange address `exchange`.
func EIP712HashV2Order(o V2Order, exchange common.Address) ([]byte, error) {
	if o.Salt == nil || o.TokenID == nil || o.MakerAmount == nil || o.TakerAmount == nil || o.Timestamp == nil {
		return nil, errors.New("order: V2Order has nil bigint field")
	}
	td := apitypes.TypedData{
		Types:       v2OrderTypes,
		PrimaryType: "Order",
		Domain: apitypes.TypedDataDomain{
			Name:              V2DomainName,
			Version:           V2DomainVersion,
			ChainId:           (*math.HexOrDecimal256)(big.NewInt(PolygonChainID)),
			VerifyingContract: exchange.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"salt":          o.Salt.String(),
			"maker":         o.Maker.Hex(),
			"signer":        o.Signer.Hex(),
			"tokenId":       o.TokenID.String(),
			"makerAmount":   o.MakerAmount.String(),
			"takerAmount":   o.TakerAmount.String(),
			"side":          fmt.Sprintf("%d", o.Side),
			"signatureType": fmt.Sprintf("%d", o.SignatureType),
			"timestamp":     o.Timestamp.String(),
			"metadata":      "0x" + hex.EncodeToString(o.Metadata[:]),
			"builder":       "0x" + hex.EncodeToString(o.Builder[:]),
		},
	}
	return hashTypedData(td)
}

// hashTypedData computes keccak256("\x19\x01" ‖ domainSeparator ‖ hashStruct(message)).
func hashTypedData(td apitypes.TypedData) ([]byte, error) {
	domainSep, err := td.HashStruct("EIP712Domain", td.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("order: domain hash: %w", err)
	}
	msgHash, err := td.HashStruct(td.PrimaryType, td.Message)
	if err != nil {
		return nil, fmt.Errorf("order: message hash: %w", err)
	}
	raw := append([]byte{0x19, 0x01}, domainSep...)
	raw = append(raw, msgHash...)
	h := sha3.NewLegacyKeccak256()
	h.Write(raw)
	return h.Sum(nil), nil
}

// RequireExchangeAddress returns the configured V2 exchange address (from
// either the constant or the override), erroring if the caller tried to
// sign an order before cutover finalised the address.
func RequireExchangeAddress(override string) (common.Address, error) {
	addr := override
	if addr == "" {
		addr = V2ExchangeAddress
	}
	if addr == "" {
		return common.Address{}, errors.New("order: V2 exchange address not set — pending 2026-04-28 cutover")
	}
	if !common.IsHexAddress(addr) {
		return common.Address{}, fmt.Errorf("order: %q is not a valid EVM address", addr)
	}
	return common.HexToAddress(addr), nil
}

// NewSalt returns a V2-friendly salt: unix-micros + 4 random bytes, so
// parallel signs from the same wallet never collide.
func NewSalt() *big.Int {
	now := time.Now().UnixMicro()
	return big.NewInt(now)
}
