package order

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Hardhat/Anvil account 0 — BIP-39 test vector that all EVM tooling agrees
// on. Any drift in mnemonic→key derivation blows up this test.
const (
	testMnemonic = "test test test test test test test test test test test junk"
	testAddress0 = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	testAddress1 = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
)

func TestWallet_KnownVectorAddress(t *testing.T) {
	w, err := NewWalletFromMnemonic(testMnemonic, "")
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	if got := w.Address().Hex(); got != testAddress0 {
		t.Fatalf("addr: got %s want %s", got, testAddress0)
	}
}

func TestWallet_CustomHDPath(t *testing.T) {
	w, err := NewWalletFromMnemonic(testMnemonic, "m/44'/60'/0'/0/1")
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	if got := w.Address().Hex(); got != testAddress1 {
		t.Fatalf("addr: got %s want %s", got, testAddress1)
	}
}

func TestWallet_EmptyMnemonicRejected(t *testing.T) {
	if _, err := NewWalletFromMnemonic("   ", ""); err == nil {
		t.Fatal("expected error for empty mnemonic")
	}
}

func TestWallet_BadMnemonicRejected(t *testing.T) {
	if _, err := NewWalletFromMnemonic("not a real mnemonic at all just words", ""); err == nil {
		t.Fatal("expected error for invalid mnemonic")
	}
}

func TestWallet_SignDigestRoundTrip(t *testing.T) {
	w, err := NewWalletFromMnemonic(testMnemonic, "")
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	// 32 random-looking bytes — doesn't need to be a real EIP-712 digest
	// for the recovery round-trip test.
	digest := crypto.Keccak256([]byte("polymarket-go phase 3 test"))
	sig, err := w.SignDigest(digest)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if len(sig) != 65 {
		t.Fatalf("sig len: got %d want 65", len(sig))
	}
	if sig[64] != 27 && sig[64] != 28 {
		t.Fatalf("sig V: got %d want 27 or 28", sig[64])
	}
	// Round-trip: recover signer from digest + signature.
	vAdj := make([]byte, 65)
	copy(vAdj, sig)
	vAdj[64] -= 27
	pub, err := crypto.SigToPub(digest, vAdj)
	if err != nil {
		t.Fatalf("recover: %v", err)
	}
	recovered := crypto.PubkeyToAddress(*pub)
	if recovered != w.Address() {
		t.Fatalf("recovered addr mismatch: got %s want %s", recovered.Hex(), w.Address().Hex())
	}
}

func TestWallet_SignDigestRejectsShort(t *testing.T) {
	w, err := NewWalletFromMnemonic(testMnemonic, "")
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	if _, err := w.SignDigest(make([]byte, 20)); err == nil {
		t.Fatal("expected error for short digest")
	}
}

func TestEIP712HashV2Order_RequiresBigInts(t *testing.T) {
	order := V2Order{
		Maker:         common.HexToAddress(testAddress0),
		Signer:        common.HexToAddress(testAddress0),
		Side:          SigBuy,
		SignatureType: SigTypeEOA,
		// salt + numeric amounts + timestamp intentionally nil
	}
	if _, err := EIP712HashV2Order(order, common.HexToAddress("0x0000000000000000000000000000000000000001")); err == nil {
		t.Fatal("expected error for nil bigint fields")
	}
}

func TestEIP712HashV2Order_Deterministic(t *testing.T) {
	order := V2Order{
		Salt:          big.NewInt(12345),
		Maker:         common.HexToAddress(testAddress0),
		Signer:        common.HexToAddress(testAddress0),
		TokenID:       big.NewInt(98765),
		MakerAmount:   big.NewInt(5_000_000), // 5 USDC (6 dec)
		TakerAmount:   big.NewInt(10_000_000),
		Side:          SigBuy,
		SignatureType: SigTypeEOA,
		Timestamp:     big.NewInt(1_700_000_000),
	}
	addr := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	h1, err := EIP712HashV2Order(order, addr)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	h2, err := EIP712HashV2Order(order, addr)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if string(h1) != string(h2) {
		t.Fatal("hash not deterministic")
	}
	if len(h1) != 32 {
		t.Fatalf("hash len: got %d want 32", len(h1))
	}
}

func TestEIP712HashV2Order_SignAndRecover(t *testing.T) {
	w, err := NewWalletFromMnemonic(testMnemonic, "")
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	order := V2Order{
		Salt:          big.NewInt(1),
		Maker:         w.Address(),
		Signer:        w.Address(),
		TokenID:       big.NewInt(42),
		MakerAmount:   big.NewInt(1_000_000),
		TakerAmount:   big.NewInt(2_000_000),
		Side:          SigBuy,
		SignatureType: SigTypeEOA,
		Timestamp:     big.NewInt(1_700_000_000),
	}
	exchange := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	digest, err := EIP712HashV2Order(order, exchange)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	sig, err := w.SignDigest(digest)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	// Recover.
	sigCopy := make([]byte, 65)
	copy(sigCopy, sig)
	sigCopy[64] -= 27
	pub, err := crypto.SigToPub(digest, sigCopy)
	if err != nil {
		t.Fatalf("recover: %v", err)
	}
	if got := crypto.PubkeyToAddress(*pub); got != w.Address() {
		t.Fatalf("recover mismatch: got %s want %s", got.Hex(), w.Address().Hex())
	}
}

func TestRequireExchangeAddress_DefaultV2(t *testing.T) {
	addr, err := RequireExchangeAddress("")
	if err != nil {
		t.Fatalf("expected V2 address to be set: %v", err)
	}
	if addr.Hex() != common.HexToAddress(V2ExchangeAddress).Hex() {
		t.Fatalf("got %s, want %s", addr.Hex(), V2ExchangeAddress)
	}
}

func TestRequireExchangeAddress_OverrideOK(t *testing.T) {
	addr, err := RequireExchangeAddress("0x1234567890abcdef1234567890abcdef12345678")
	if err != nil {
		t.Fatalf("override: %v", err)
	}
	if !strings.EqualFold(addr.Hex(), "0x1234567890AbcdEF1234567890aBcdef12345678") {
		t.Fatalf("addr: got %s", addr.Hex())
	}
}

func TestRequireExchangeAddress_BadHex(t *testing.T) {
	if _, err := RequireExchangeAddress("not-an-address"); err == nil {
		t.Fatal("expected error for bad hex")
	}
}
