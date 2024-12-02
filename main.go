package main

import (
	"log"
	"sync"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
	"github.com/SunTzu71/suntzu_blockchain/constants"
)

// Inialize blockchain name
func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ": ")
}

// Main function to run the blockchain
func main() {

	var wg sync.WaitGroup

	genesisBlock := blockchain.NewBlock("0x0", 0)
	transaction := blockchain.NewTransaction("0x0", "0x1", 2000, []byte("This is a test transaction"))
	blockchain := blockchain.NewBlockchain(*genesisBlock)
	log.Println(blockchain.ToJson())
	log.Println("Start Mining")
	wg.Add(1)
	go blockchain.ProofOfWorkMining("SunTzu")
	blockchain.AddTransactionToTransactionPool(*transaction)
	wg.Wait()
}
