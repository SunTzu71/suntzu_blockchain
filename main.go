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
	block := blockchain.NewBlock("0x", 0)
	log.Println(block)

	transaction := blockchain.NewTransaction("0x0", "0x1", 2000, []byte("This is a test transaction"))
	log.Println(transaction)
}
