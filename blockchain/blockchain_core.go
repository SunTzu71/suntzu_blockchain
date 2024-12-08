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
	Address         string         `json:"address"`
}

// NewBlockchain: creates a new blockchain instance with a genesis block
// If blockchain data exists in the database (checked via DBKeyExists), retrieves and returns it
// Otherwise creates a new blockchain with the genesis block and persists it via DBAddBlockchain
// Returns a pointer to the BlockchainCore instance in either case
func NewBlockchain(genesisBlock Block, address string) *BlockchainCore {
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
		blockchainCore.Address = address

		err := DBAddBlockchain(*blockchainCore)
		if err != nil {
			log.Fatal(err)
		}

		return blockchainCore
	}
}

// NewBlockchainSync: creates a copy of an existing blockchain with a new address
// Takes a pointer to an existing BlockchainCore and a new address string
// Returns a pointer to the new BlockchainCore instance with updated address
func NewBlockchainSync(bc1 *BlockchainCore, address string) *BlockchainCore {
	bc2 := bc1
	bc2.Address = address

	err := DBAddBlockchain(*bc2)
	if err != nil {
		log.Fatal(err)
	}
	return bc2
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

// AddTransactionToTransactionPool: verifies a transaction's signature and sender balance,
// sets its status accordingly, and adds it to the blockchain's transaction pool.
// The transaction is verified for valid signature and sufficient sender balance.
// Updates transaction status to success or failure based on verification.
// Persists the updated blockchain to database after adding the transaction.
func (bc *BlockchainCore) AddTransactionToTransactionPool(transaction *Transaction) {
	validTransaction := transaction.VerifyTransaction()

	realBalance := bc.CalculateTotalCrypto(transaction.From)
	validRealBalance := realBalance >= transaction.Value

	if validTransaction && validRealBalance {
		transaction.Status = constants.TRANSACTION_VERIFY_SUCCESS
	} else {
		transaction.Status = constants.TRANSACTION_VERIFY_FAILED
	}
	transaction.PublicKey = ""
	bc.TransactionPool = append(bc.TransactionPool, transaction)

	// Save the blockchain to the database
	err := DBAddBlockchain(*bc)
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
	err := DBAddBlockchain(*bc)
	if err != nil {
		log.Fatal(err)
	}
}

// ProofOfWorkMining continuously mines new blocks using proof of work consensus.
// It takes a miner's address as input and rewards successful mining with coins.
// The function runs indefinitely, creating new blocks that meet the mining difficulty
// requirement by incrementing a nonce value until a valid hash is found.
func (bc *BlockchainCore) ProofOfWorkMining(minersAddress string) {
	log.Println("Proof of work mining started")
	// calculcate the prevHash
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash()

	// had to set this as int64 getting error that new block nonce was int
	var nonce int64 = 0

	for {
		guessBlock := NewBlock(prevHash, nonce)

		for _, txn := range bc.TransactionPool {
			newTxn := new(Transaction)
			newTxn.Data = txn.Data
			newTxn.From = txn.From
			newTxn.To = txn.To
			newTxn.Status = txn.Status
			newTxn.Timestamp = txn.Timestamp
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

			prevHash = bc.Blocks[len(bc.Blocks)-1].Hash()
			nonce = 0
			continue
		}
		nonce++
	}
}

// CalculateTotalCrypto: calculates the total balance of cryptocurrency for a given address
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

// GetAllTransactions: retrieves all transactions from both the transaction pool and blocks
// in reverse chronological order (newest first). It first collects transactions from the
// transaction pool, then adds transactions from blocks, excluding mining reward transactions
// (those from BLOCKCHAIN_ADDRESS). Returns a slice of all non-reward transactions.
func (bc *BlockchainCore) GetAllNonRewardedTransactions() []Transaction {

	newestTxns := []Transaction{}

	for i := len(bc.TransactionPool) - 1; i >= 0; i-- {
		newestTxns = append(newestTxns, *bc.TransactionPool[i])
	}

	txns := []Transaction{}

	for _, block := range bc.Blocks {
		for _, txn := range block.Transactions {
			if txn.From != constants.BLOCKCHAIN_ADDRESS {
				txns = append(txns, *txn)
			}
		}
	}

	for i := len(txns) - 1; i >= 0; i-- {
		newestTxns = append(newestTxns, txns[i])
	}

	return newestTxns
}
