package main

import (
	"log"
	"sync"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
	"github.com/SunTzu71/suntzu_blockchain/constants"
	"github.com/SunTzu71/suntzu_blockchain/wallet"
)

// Inialize blockchain name
func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ": ")
}

// Main function to run the blockchain
func main() {

	var wg sync.WaitGroup
	wallet2, _ := wallet.NewWallet()

	genesisBlock := blockchain.NewBlock("0x0", 0)
	blockchain1 := blockchain.NewBlockchain(*genesisBlock)

	// wallet1, _ := wallet.NewWallet()
	// uTxn := blockchain.NewTransaction(wallet1.GetAddress(), wallet2.GetAddress(), 1000, []byte("This is a test transaction"))
	// sTxn, _ := wallet1.GetSignedTransaction(*uTxn)
	// for i := 0; i < 10; i++ {
	// 	blockchain1.AddTransactionToTransactionPool(*sTxn)
	// }

	log.Println(blockchain1.ToJson())
	log.Println("Start Mining")
	wg.Add(1)
	go blockchain1.ProofOfWorkMining(wallet2.GetAddress())
	wg.Wait()

}
