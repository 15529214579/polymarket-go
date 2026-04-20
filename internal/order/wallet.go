package order

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// DefaultHDPath is BIP-44 m/44'/60'/0'/0/0 — MetaMask / Anvil / Hardhat
// account 0. Polymarket-Go wallet was derived at this path and stored in
// Bitwarden (see TOOLS.md + PRINCIPLES.md P4).
const DefaultHDPath = "m/44'/60'/0'/0/0"

// Wallet wraps an in-memory ECDSA key derived from a BIP-39 mnemonic.
// The private key is never written to disk or logs. The mnemonic is
// consumed once at construction and then discarded by the caller.
type Wallet struct {
	privKey *ecdsa.PrivateKey
	address common.Address
}

// NewWalletFromMnemonic derives an EVM key from a BIP-39 mnemonic at the
// given HD path (empty → DefaultHDPath).
func NewWalletFromMnemonic(mnemonic, hdPath string) (*Wallet, error) {
	mnemonic = strings.TrimSpace(mnemonic)
	if mnemonic == "" {
		return nil, errors.New("order: empty mnemonic")
	}
	if hdPath == "" {
		hdPath = DefaultHDPath
	}
	w, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("order: bad mnemonic: %w", err)
	}
	path, err := hdwallet.ParseDerivationPath(hdPath)
	if err != nil {
		return nil, fmt.Errorf("order: bad hd path %q: %w", hdPath, err)
	}
	acct, err := w.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("order: derive: %w", err)
	}
	pk, err := w.PrivateKey(acct)
	if err != nil {
		return nil, fmt.Errorf("order: extract key: %w", err)
	}
	return &Wallet{privKey: pk, address: acct.Address}, nil
}

// Address returns the checksummed EVM address corresponding to the wallet.
func (w *Wallet) Address() common.Address { return w.address }

// SignDigest signs a 32-byte EIP-712 digest and returns a 65-byte signature
// with V ∈ {27, 28}, matching the encoding accepted by Polymarket CLOB and
// most on-chain contracts.
func (w *Wallet) SignDigest(digest []byte) ([]byte, error) {
	if len(digest) != 32 {
		return nil, fmt.Errorf("order: digest must be 32 bytes, got %d", len(digest))
	}
	sig, err := crypto.Sign(digest, w.privKey)
	if err != nil {
		return nil, err
	}
	// go-ethereum returns V in {0,1}; EIP-712 callers expect {27,28}.
	if sig[64] < 27 {
		sig[64] += 27
	}
	return sig, nil
}

// LoadMnemonicFromBitwarden shells out to the `bw` CLI to read a secure
// note and return a custom field by name. Intended to run once at daemon
// start; the returned string should be handed directly to
// NewWalletFromMnemonic and dropped. Never logged, never written to disk.
//
// Pre-req: BW_SESSION must be set in the environment (see TOOLS.md).
func LoadMnemonicFromBitwarden(itemName, fieldName string) (string, error) {
	itemName = strings.TrimSpace(itemName)
	fieldName = strings.TrimSpace(fieldName)
	if itemName == "" || fieldName == "" {
		return "", errors.New("order: bitwarden itemName and fieldName required")
	}
	// Call `bw` directly — no `bash -c`, no shell interpolation — and
	// parse the JSON locally. Side-steps shell-injection and gosec G204.
	out, err := exec.Command("bw", "get", "item", itemName).Output()
	if err != nil {
		return "", fmt.Errorf("order: bw get item %q: %w", itemName, err)
	}
	var item struct {
		Fields []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(out, &item); err != nil {
		return "", fmt.Errorf("order: parse bw response: %w", err)
	}
	for _, f := range item.Fields {
		if f.Name == fieldName {
			v := strings.TrimSpace(f.Value)
			if v == "" {
				return "", fmt.Errorf("order: bitwarden field %q on item %q is empty", fieldName, itemName)
			}
			return v, nil
		}
	}
	return "", fmt.Errorf("order: bitwarden field %q not found on item %q", fieldName, itemName)
}
