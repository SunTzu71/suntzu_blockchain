package blockchain

import "encoding/json"

type BlockchainCore struct {
	TransactionPool []*Transaction `json:"transaction_pool"`
	Blocks          []*Block       `json:"blocks"`
}

// NewBlockchain creates a new blockchain instance with a genesis block
// Returns a pointer to the new BlockchainCore
func NewBlockchain(genesisBlock Block) *BlockchainCore {
	blockchainCore := new(BlockchainCore)
	blockchainCore.TransactionPool = []*Transaction{}
	blockchainCore.Blocks = []*Block{}
	blockchainCore.Blocks = append(blockchainCore.Blocks, &genesisBlock)

	return blockchainCore
}

// ToJson converts the BlockchainCore structure to a JSON string
// Returns the JSON string representation or an error message if marshal fails
func (bc BlockchainCore) ToJson() string {
	nb, err := json.Marshal(bc)
	if err != nil {
		return err.Error()
	}

	return string(nb)
}
