package blockchain

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/SunTzu71/suntzu_blockchain/constants"
)

type BlockchainCore struct {
	TransactionPool []*Transaction `json:"transaction_pool"`
	Blocks          []*Block       `json:"blocks"`
}

// NewBlockchain creates a new blockchain instance with a genesis block
// Uses KeyExists() to check if blockchain data exists in database
// Uses PutIntoDb() to persist blockchain data
// Returns a pointer to the new BlockchainCore
func NewBlockchain(genesisBlock Block) *BlockchainCore {
	if DBKeyExists() {
		blockchianCore, err := DBGetBlockchain()
		if err != nil {
			log.Fatal(err)
		}

		return blockchianCore
	} else {
		blockchainCore := new(BlockchainCore)
		blockchainCore.TransactionPool = []*Transaction{}
		blockchainCore.Blocks = []*Block{}
		blockchainCore.Blocks = append(blockchainCore.Blocks, &genesisBlock)

		err := DBAddBllockchain(*blockchainCore)
		if err != nil {
			log.Fatal(err)
		}

		return blockchainCore
	}
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

// AddTransactionToTransactionPool takes a Transaction and adds it to the blockchain's transaction pool
func (bc *BlockchainCore) AddTransactionToTransactionPool(transaction Transaction) {
	bc.TransactionPool = append(bc.TransactionPool, &transaction)

	// Save the blockchain to the database
	err := DBAddBllockchain(*bc)
	if err != nil {
		log.Fatal(err)
	}
}

// AddBlock adds a new block to the blockchain and removes its transactions from the transaction pool.
// It takes a pointer to a Block as input and updates both the blockchain's transaction pool
// and blocks array. Transactions in the new block are removed from the pool to prevent double-spending.
func (bc *BlockchainCore) AddBlock(b *Block) {
	// Create a map of transaction hashes in the new block
	txnMap := make(map[string]bool)
	for _, txn := range b.Transactions {
		txnMap[txn.TransactionHash] = true
	}

	// Create a new slice for transactions that should remain in the pool
	var newTransactionPool []*Transaction
	for _, txn := range bc.TransactionPool {
		if !txnMap[txn.TransactionHash] {
			newTransactionPool = append(newTransactionPool, txn)
		}
	}

	// Replace the transaction pool with the filtered version
	bc.TransactionPool = newTransactionPool

	// Add block to blockchain
	bc.Blocks = append(bc.Blocks, b)

	// Save the blockchain to the database
	err := DBAddBllockchain(*bc)
	if err != nil {
		log.Fatal(err)
	}
}

// ProofOfWorkMining continuously mines new blocks using proof of work consensus.
// It takes a miner's address as input and rewards successful mining with coins.
// The function runs indefinitely, creating new blocks that meet the mining difficulty
// requirement by incrementing a nonce value until a valid hash is found.
func (bc *BlockchainCore) ProofOfWorkMining(minersAddress string) {
	// calculcate the prevHash
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash()

	// had to set this as int64 getting error that new block nonce was int
	var nonce int64 = 0

	for {
		// create a new block
		guessBlock := NewBlock(prevHash, nonce)

		// copy the transaction pool
		for _, txn := range bc.TransactionPool {
			//newTxn := NewTransaction(txn.From, txn.To, txn.Value, txn.Data)
			newTxn := new(Transaction)
			newTxn.Data = txn.Data
			newTxn.From = txn.From
			newTxn.To = txn.To
			newTxn.Status = txn.Status
			newTxn.Value = txn.Value
			newTxn.TransactionHash = txn.TransactionHash
			newTxn.PublicKey = txn.PublicKey
			newTxn.Signature = txn.Signature

			guessBlock.AddTransactionToTheBlock(newTxn)
		}

		// guess the hash
		guessHash := guessBlock.Hash()
		desiredHash := strings.Repeat("0", constants.MINING_DIFFICULTY)
		ourSolutionHash := guessHash[2 : 2+constants.MINING_DIFFICULTY]

		if ourSolutionHash == desiredHash {
			rewardTxn := NewTransaction(constants.BLOCKCHAIN_ADDRESS, minersAddress, constants.MINING_REWARD, []byte{})
			rewardTxn.Status = constants.SUCCESS
			guessBlock.Transactions = append(guessBlock.Transactions, rewardTxn)
			bc.AddBlock(guessBlock)

			log.Println(bc.ToJson())

			prevHash = bc.Blocks[len(bc.Blocks)-1].Hash()
			nonce = 0
			continue
		}
		nonce++
	}
}

// CalculateTotalCrypto calculates the total balance of cryptocurrency for a given address
// by examining all successful transactions in both the blockchain and transaction pool.
// It adds received amounts (To) and subtracts sent amounts (From) for the address.
func (bc *BlockchainCore) CalculateTotalCrypto(address string) uint64 {
	var balance uint64 = 0

	for _, block := range bc.Blocks {
		for _, txn := range block.Transactions {
			if txn.Status == constants.SUCCESS {
				if txn.To == address {
					balance += txn.Value
				}
				if txn.From == address {
					balance -= txn.Value
				}
			}
		}
	}

	for _, txn := range bc.TransactionPool {
		if txn.Status == constants.SUCCESS {
			if txn.To == address {
				balance += txn.Value
			}
			if txn.From == address {
				balance -= txn.Value
			}
		}
	}

	return balance
}
