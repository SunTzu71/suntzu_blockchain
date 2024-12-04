package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"math/big"
	"time"

	"github.com/SunTzu71/suntzu_blockchain/constants"
)

type Transaction struct {
	From            string `json:"from"`
	To              string `json:"to"`
	Value           uint64 `json:"value"`
	Data            []byte `json:"data"`
	Status          string `json:"status"`
	Timestamp       uint64 `json:"timestamp"`
	TransactionHash string `json:"transaction_hash"`
	PublicKey       string `json:"public_key,omitempty"`
	Signature       []byte `json:"Signature"`
}

// NewTransaction creates and returns a new Transaction object initialized with the provided parameters
// and default values for Status, PublicKey, and Signature fields
func NewTransaction(from string, to string, value uint64, data []byte) *Transaction {
	t := new(Transaction)
	t.From = from
	t.To = to
	t.Value = value
	t.Data = data
	t.Status = constants.PENDING
	t.Timestamp = uint64(time.Now().Unix())
	t.TransactionHash = t.Hash()
	t.PublicKey = ""
	t.Signature = []byte{}
	return t
}

// ToJson converts a Transaction object to its JSON string representation
func (t Transaction) ToJson() string {
	nb, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(nb)
}

// VerifyTransaction checks if the transaction is valid by verifying:
// 1. The value is not zero
// 2. The value does not exceed maximum uint64
// 3. The signature is valid
// Returns true if all checks pass, false otherwise
func (t Transaction) VerifyTransaction() bool {
	if t.Value == 0 {
		return false
	}

	if t.Value > math.MaxUint64 {
		return false
	}

	valid := t.VeryifySignature()
	if !valid {
		return false
	}

	return true
}

// VeryifySignature verifies the digital signature of a transaction using ECDSA
// It first checks if signature and public key exist, then verifies the signature
// against the transaction hash using the public key.
// Returns true if signature is valid, false otherwise
func (t Transaction) VeryifySignature() bool {
	if t.Signature == nil || t.PublicKey == "" {
		return false
	}

	signature := t.Signature
	publicKeyHex := t.PublicKey
	t.Signature = []byte{}
	t.PublicKey = ""
	publicKeyEcdsa := GetPublicKeyFromHex(publicKeyHex)

	bs, _ := json.Marshal(t)
	hash := sha256.Sum256(bs)

	valid := ecdsa.VerifyASN1(publicKeyEcdsa, hash[:], signature)
	t.Signature = signature

	return valid
}

// Hash generates a SHA-256 hash of the transaction data and returns it as a hex string with prefix.
// The transaction is first marshaled to JSON, then hashed, and finally encoded to a hex string.
func (t Transaction) Hash() string {
	bs, _ := json.Marshal(t)
	sum := sha256.Sum256(bs)
	hexRep := hex.EncodeToString(sum[:32])
	formattedHexRep := constants.HEX_PREFIX + hexRep

	return formattedHexRep
}

// GetPublicKeyFromHex converts a hex string representation of a public key to an ECDSA public key
// It strips the hex prefix, splits the remaining string into x and y coordinates,
// and creates a new public key using the P256 curve
func GetPublicKeyFromHex(publicKeyHex string) *ecdsa.PublicKey {
	rpk := publicKeyHex[2:]
	xHex := rpk[0:64]
	yHex := rpk[64:]
	x := new(big.Int)
	y := new(big.Int)
	x.SetString(xHex, 16)
	y.SetString(yHex, 16)

	var npk ecdsa.PublicKey
	npk.Curve = elliptic.P256()
	npk.X = x
	npk.Y = y

	return &npk
}
