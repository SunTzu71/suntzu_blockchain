package main

import (
	"log"

	"github.com/SunTzu71/suntzu_blockchain/constants"
)

// Inialize blockchain name
func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ": ")
}

// Main function to run the blockchain
func main() {

	// var wg sync.WaitGroup

	// genesisBlock := blockchain.NewBlock("0x0", 0)
	// transaction := blockchain.NewTransaction("0x0", "0x1", 2000, []byte("This is a test transaction"))
	// blockchain := blockchain.NewBlockchain(*genesisBlock)
	// log.Println(blockchain.ToJson())
	// log.Println("Start Mining")
	// wg.Add(1)
	// go blockchain.ProofOfWorkMining("SunTzu")
	// blockchain.AddTransactionToTransactionPool(*transaction)
	// wg.Wait()

	// wallet1, _ := wallet.NewWallet()
	// log.Println("private key: ", wallet1.GetPrivateKeyHex())
	// log.Println("public key: ", wallet1.GetPublicKeyHex())
	// log.Println("address: ", wallet1.GetAddress())

	// wallet2 := wallet.NewWalletFromPrivateKeyHex(wallet1.GetPrivateKeyHex())
	// log.Println("private key: ", wallet2.GetPrivateKeyHex())
	// log.Println("public key: ", wallet2.GetPublicKeyHex())
	// log.Println("address: ", wallet2.GetAddress())

	// log.Println("Checking Equals--------")
	// log.Println("private key: ", wallet1.GetPrivateKeyHex() == wallet2.GetPrivateKeyHex())
	// log.Println("public key: ", wallet1.GetPublicKeyHex() == wallet2.GetPublicKeyHex())
	// log.Println("address: ", wallet1.GetAddress() == wallet2.GetAddress())
}
