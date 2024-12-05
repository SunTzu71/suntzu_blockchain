package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
	"github.com/SunTzu71/suntzu_blockchain/blockchainserver"
	"github.com/SunTzu71/suntzu_blockchain/constants"
	"github.com/SunTzu71/suntzu_blockchain/walletserver"
)

// Inialize blockchain name
func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ": ")
}

// Main function to run the blockchain
func main() {

	genesisBlock := blockchain.NewBlock("0x0", 0)
	blockchain1 := blockchain.NewBlockchain(*genesisBlock)
	bcs := blockchainserver.CreateBlockchainServer(8000, blockchain1)
	ws := walletserver.CreateWalletServer(8080, "http://127.0.0.1:8000")

	// Start the blockchain server and wallet server
	go bcs.StartBlockchainServer()
	go ws.StartWalletServer()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
