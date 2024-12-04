package main

import (
	"log"

	"github.com/SunTzu71/suntzu_blockchain/constants"
	"github.com/SunTzu71/suntzu_blockchain/walletserver"
)

// Inialize blockchain name
func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ": ")
}

// Main function to run the blockchain
func main() {
	// genesisBlock := blockchain.NewBlock("0x0", 0)
	// blockchain1 := blockchain.NewBlockchain(*genesisBlock)
	// bcs := blockchainserver.CreateBlockchainServer(8000, blockchain1)
	// bcs.StartBlockchainServer()

	ws := walletserver.CreateWalletServer(8080, "abcd")
	ws.StartWalletServer()
}
