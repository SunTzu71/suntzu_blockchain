package blockchain

import "encoding/json"

type BlockchainCore struct {
	TransactionPool []*Transaction `json:"transaction_pool"`
	Blocks          []*Block       `json:"blocks"`
}

func NewBlockchain(genesisBlock Block) *BlockchainCore {
	blockchainCore := new(BlockchainCore)
	blockchainCore.TransactionPool = []*Transaction{}
	blockchainCore.Blocks = []*Block{}
	blockchainCore.Blocks = append(blockchainCore.Blocks, &genesisBlock)

	return blockchainCore
}

func (bc BlockchainCore) ToJson() string {
	nb, err := json.Marshal(bc)
	if err != nil {
		return err.Error()
	}

	return string(nb)
}
