package main

import (
	"log"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
	"github.com/SunTzu71/suntzu_blockchain/constants"
)

// Inialize blockchain name
func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ": ")
}

// Main function to run the blockchain
func main() {

	genesisBlock := blockchain.NewBlock("0x0", 0)
	transaction := blockchain.NewTransaction("0x0", "0x1", 2000, []byte("This is a test transaction"))
	genesisBlock.Transactions = append(genesisBlock.Transactions, transaction)
	blockchain := blockchain.NewBlockchain(*genesisBlock)
	log.Println(blockchain.ToJson())
}
