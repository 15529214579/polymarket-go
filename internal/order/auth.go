package order

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const (
	clobAuthDomainName    = "ClobAuthDomain"
	clobAuthDomainVersion = "1"
)

type APICredentials struct {
	APIKey     string `json:"apiKey"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

var clobAuthTypes = apitypes.Types{
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
	},
	"ClobAuth": []apitypes.Type{
		{Name: "address", Type: "address"},
		{Name: "timestamp", Type: "string"},
		{Name: "nonce", Type: "uint256"},
		{Name: "message", Type: "string"},
	},
}

func buildClobAuthDigest(addr common.Address, ts int64, nonce int) ([]byte, error) {
	td := apitypes.TypedData{
		Types:       clobAuthTypes,
		PrimaryType: "ClobAuth",
		Domain: apitypes.TypedDataDomain{
			Name:    clobAuthDomainName,
			Version: clobAuthDomainVersion,
			ChainId: (*math.HexOrDecimal256)(big.NewInt(PolygonChainID)),
		},
		Message: apitypes.TypedDataMessage{
			"address":   addr.Hex(),
			"timestamp": strconv.FormatInt(ts, 10),
			"nonce":     fmt.Sprintf("%d", nonce),
			"message":   "This message attests that I control the given wallet",
		},
	}
	return hashTypedData(td)
}

func DeriveAPIKey(clobBase string, w *Wallet) (*APICredentials, error) {
	ts := time.Now().Unix()
	nonce := 0

	digest, err := buildClobAuthDigest(w.Address(), ts, nonce)
	if err != nil {
		return nil, fmt.Errorf("order: clob auth digest: %w", err)
	}
	sig, err := w.SignDigest(digest)
	if err != nil {
		return nil, fmt.Errorf("order: clob auth sign: %w", err)
	}

	req, err := http.NewRequest("GET", clobBase+"/auth/derive-api-key", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("POLY_ADDRESS", w.Address().Hex())
	req.Header.Set("POLY_SIGNATURE", "0x"+fmt.Sprintf("%x", sig))
	req.Header.Set("POLY_TIMESTAMP", strconv.FormatInt(ts, 10))
	req.Header.Set("POLY_NONCE", strconv.Itoa(nonce))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("order: derive api key: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("order: derive api key %d: %s", resp.StatusCode, body)
	}

	var creds APICredentials
	if err := json.Unmarshal(body, &creds); err != nil {
		return nil, fmt.Errorf("order: parse api creds: %w", err)
	}
	if creds.APIKey == "" || creds.Secret == "" {
		return nil, fmt.Errorf("order: empty api credentials in response: %s", body)
	}
	return &creds, nil
}

func buildL2Headers(creds *APICredentials, addr common.Address, method, path, bodyJSON string) http.Header {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	msg := ts + strings.ToUpper(method) + path + bodyJSON
	secretBytes, _ := base64.URLEncoding.DecodeString(creds.Secret)
	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(msg))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	h := http.Header{}
	h.Set("POLY_ADDRESS", addr.Hex())
	h.Set("POLY_SIGNATURE", sig)
	h.Set("POLY_TIMESTAMP", ts)
	h.Set("POLY_API_KEY", creds.APIKey)
	h.Set("POLY_PASSPHRASE", creds.Passphrase)
	h.Set("Content-Type", "application/json")
	return h
}
