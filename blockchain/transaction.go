package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"

	"github.com/SunTzu71/suntzu_blockchain/constants"
)

type Transaction struct {
	From            string `json:"from"`
	To              string `json:"to"`
	Value           uint64 `json:"value"`
	Data            []byte `json:"data"`
	Status          string `json:"status"`
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

// VerifyTxn validates a transaction by checking if the value is within valid range.
// Returns false if value is 0 or exceeds uint64 max, true otherwise.
func (t Transaction) VerifyTransaction() bool {
	if t.Value == 0 {
		return false
	}

	if t.Value > math.MaxUint64 {
		return false
	}

	// TODO: need to implement VerifySignature when we create wallet logic
	return true
}

func (t Transaction) Hash() string {
	bs, _ := json.Marshal(t)
	sum := sha256.Sum256(bs)
	hexRep := hex.EncodeToString(sum[:32])
	formattedHexRep := constants.HEX_PREFIX + hexRep

	return formattedHexRep
}
