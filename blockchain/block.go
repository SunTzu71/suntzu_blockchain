package blockchain

import (
	"encoding/json"
	"time"
)

type Block struct {
	PrevHash     string         `json:"prev_hash"`
	Timestamp    int64          `json:"timestamp"`
	Nonce        int64          `json:"nonce"`
	Transactions []*Transaction `json:"transactions"`
}

// NewBlock creates a new Block instance with the provided previous hash and nonce value,
// initializing its timestamp to the current time and an empty transaction list
func NewBlock(prevHash string, nonce int64) *Block {
	block := new(Block)
	block.PrevHash = prevHash
	block.Timestamp = time.Now().UnixNano()
	block.Nonce = nonce
	block.Transactions = []*Transaction{}

	return block
}

// ToJson converts a Block instance into a JSON string representation.
// If there's an error during marshaling, it returns the error message as a string.
func (b Block) ToJson() string {
	nb, err := json.Marshal(b)
	if err != nil {
		return err.Error()
	}

	return string(nb)
}
