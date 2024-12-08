package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

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
	chainCommandSet := flag.NewFlagSet("chain", flag.ExitOnError)
	walletCommandSet := flag.NewFlagSet("wallet", flag.ExitOnError)

	chainPort := chainCommandSet.Uint("port", 8000, "port to run the blockchain server")
	chainMiner := chainCommandSet.String("miner", "", "miner address")
	remoteNode := chainCommandSet.String("remote_node", "", "remote node address")

	walletPort := walletCommandSet.Uint("port", 8080, "port to run the wallet server")
	blockchainNodeAddress := walletCommandSet.String("node", "http://127.0.0.1:8000", "blockchain node address")

	if len(os.Args) < 2 {
		fmt.Println("Error: expected chain or wallet command")
	}

	switch os.Args[1] {
	case "chain":
		chainCommandSet.Parse(os.Args[2:])
		if chainCommandSet.Parsed() {

			if *chainMiner == "" || chainCommandSet.NFlag() == 0 {
				fmt.Println("Usage of chain subcommand: ")
				chainCommandSet.PrintDefaults()
				os.Exit(1)
			}

			// if remote node is empty launch new blockchain
			if *remoteNode == "" {
				genesisBlock := blockchain.NewBlock("0x0", 0)
				blockchain := blockchain.NewBlockchain(*genesisBlock, "http://127.0.0.1:"+strconv.Itoa(int(*chainPort)))
				bcs := blockchainserver.CreateBlockchainServer(uint64(*chainPort), blockchain)
				go bcs.StartBlockchainServer()
				go bcs.BlockchainPtr.ProofOfWorkMining(*chainMiner)

				// Wait for interrupt signal
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				<-c
			} else {
				blockchain1, err := blockchain.SyncBlockchain(*remoteNode)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				blockchain2 := blockchain.NewBlockchainSync(blockchain1, "http://127.0.0.1:"+strconv.Itoa(int(*chainPort)))
				bcs := blockchainserver.CreateBlockchainServer(uint64(*chainPort), blockchain2)
				go bcs.StartBlockchainServer()
				go bcs.BlockchainPtr.ProofOfWorkMining(*chainMiner)

				// Wait for interrupt signal
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				<-c
			}
		}
	case "wallet":
		walletCommandSet.Parse(os.Args[2:])
		if walletCommandSet.Parsed() {
			if walletCommandSet.NFlag() == 0 {
				fmt.Println("Usage of wallet subcommand: ")
				walletCommandSet.PrintDefaults()
				os.Exit(1)
			}
			ws := walletserver.CreateWalletServer(uint16(*walletPort), *blockchainNodeAddress)
			go ws.StartWalletServer()

			// Wait for interrupt signal
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
		}
	default:
		fmt.Println("Error: expected chain or wallet command")
		os.Exit(1)
	}
}
