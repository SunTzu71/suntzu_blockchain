package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/SunTzu71/suntzu_blockchain/constants"
)

type Block struct {
	BlockNumber  uint64         `json:"block_number"`
	PrevHash     string         `json:"prev_hash"`
	Timestamp    int64          `json:"timestamp"`
	Nonce        int64          `json:"nonce"`
	Transactions []*Transaction `json:"transactions"`
}

// NewBlock creates a new Block instance with the provided previous hash and nonce value,
// initializing its timestamp to the current time and an empty transaction list
func NewBlock(prevHash string, nonce int64, blockNumber uint64) *Block {
	block := new(Block)
	block.PrevHash = prevHash
	block.Timestamp = time.Now().Unix()
	block.Nonce = nonce
	block.Transactions = []*Transaction{}
	block.BlockNumber = blockNumber

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

// Hash calculates and returns a SHA-256 hash of the block's contents as a hexadecimal string,
// prefixed with the hex prefix constant
func (b Block) Hash() string {
	bs, _ := json.Marshal(b)
	sum := sha256.Sum256(bs)
	hexRep := hex.EncodeToString(sum[:32])
	formattedHexRep := constants.HEX_PREFIX + hexRep

	return formattedHexRep
}

// AddTransactionToTheBlock: adds a transaction to the block's transaction list and updates
// its status to either SUCCESS or FAILED based on its verification status.
func (b *Block) AddTransactionToTheBlock(txn *Transaction) {

	if txn.Status == constants.TRANSACTION_VERIFY_SUCCESS {
		txn.Status = constants.SUCCESS
	} else {
		txn.Status = constants.FAILED
	}

	b.Transactions = append(b.Transactions, txn)
}
