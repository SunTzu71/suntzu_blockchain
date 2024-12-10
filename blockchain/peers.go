package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SunTzu71/suntzu_blockchain/constants"
)

// SyncBlockchain: retrieves and synchronizes the blockchain from a given address.
// It makes an HTTP GET request to the provided address, reads the blockchain data,
// and unmarshals it into a BlockchainCore struct. Returns a pointer to the
// synchronized blockchain and any errors encountered.
func SyncBlockchain(address string) (*BlockchainCore, error) {
	log.Println("Syncing blockchain from:", address)
	outURL := fmt.Sprintf("%s/", address)
	resp, err := http.Get(outURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs BlockchainCore
	err = json.Unmarshal(data, &bs)
	if err != nil {
		return nil, err
	}

	log.Println("Blockchain synced:", address)

	return &bs, nil
}

// UpdatePeers: updates the peers map in the blockchain with the provided peers map.
// Takes a map of peer addresses to boolean values. Uses mutex locking to ensure
// thread safety when updating the peers. After updating, saves the blockchain state
// to the database.
func (bc *BlockchainCore) UpdatePeers(peers map[string]bool) {
	mutex.Lock()
	defer mutex.Unlock()

	log.Println("Updating peers list...", peers)

	bc.Peers = peers

	err := DBAddBlockchain(*bc)
	if err != nil {
		log.Fatal(err)
	}
}

// SendPeersList: sends the blockchain's peer list to a specified address via HTTP POST.
// Converts the peer list to JSON and sends it to the /send-peers-list endpoint
// at the given address.
func (bc *BlockchainCore) SendPeersList(address string) {
	data := bc.PeersToJson()
	ourURL := fmt.Sprintf("%s/send-peers-list", address)
	resp, err := http.Post(ourURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error sending peers list: %v", err)
		return
	}
	defer resp.Body.Close()
}

// CheckStatus: checks if a blockchain server at the given address is available and running.
// Makes an HTTP GET request to the server's status endpoint and verifies the response matches
// the expected blockchain status value. Returns true if the server is running and accessible,
// false otherwise.
func (bc *BlockchainCore) CheckStatus(address string) bool {
	outURL := fmt.Sprintf("%s/check-server-status", address)
	resp, err := http.Get(outURL)
	if err != nil {
		log.Println("Error checking server status:", err)
		return false
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading server status response:", err)
		return false
	}

	return string(data) == constants.BLOCKCHAIN_STATUS
}

// BroadcastPeerList: broadcasts the blockchain's peer list to all active peers in the network.
// Iterates through the peer list, sending the peer list to each active peer except itself.
// Includes a delay between broadcasts to prevent network congestion.
func (bc *BlockchainCore) BroadcastPeerList() {
	for peer, status := range bc.Peers {
		if peer != bc.Address && status {
			bc.SendPeersList(peer)
			time.Sleep(constants.PEER_LIST_UPDATE_INTERVAL * time.Second)
		}
	}
}

// DialUpdatePeers: continuously checks and updates the status of peers in the blockchain network.
// Iterates through the peer list periodically, checking each peer's status via HTTP
// and updating the peers map accordingly. The blockchain's own address is always marked as active.
// After updating peers, broadcasts the new peer list to the network and sleeps for the configured
// ping interval before the next update cycle.
func (bc *BlockchainCore) DialUpdatePeers() {
	ticker := time.NewTicker(constants.PEER_PING_INTERVAL * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Pinging peers", bc.Peers)
			newList := make(map[string]bool)
			for peer := range bc.Peers {
				if peer != bc.Address {
					newList[peer] = bc.CheckStatus(peer)
					time.Sleep(constants.PEER_PING_INTERVAL * time.Second)
				} else {
					newList[peer] = true
				}
			}

			bc.UpdatePeers(newList)
			log.Println("Peers updated")

			bc.BroadcastPeerList()
		}
	}
}

// SendTransactionPeer: sends a transaction to a specified peer address via HTTP POST.
// Takes the peer's address and a transaction pointer as input, converts the transaction
// to JSON, and sends it to the peer's /send-transaction endpoint.
func (bc *BlockchainCore) SendTransactionPeer(address string, txn *Transaction) {
	data := txn.ToJson()
	ourURL := fmt.Sprintf("%s/send-transaction", address)
	http.Post(ourURL, "application/json", strings.NewReader(data))
}

// BroadcastTransaction: broadcasts a transaction to all active peers in the network.
// Iterates through the peer list, sending the transaction to each active peer except itself.
// Includes a delay between broadcasts to prevent network congestion.
func (bc *BlockchainCore) BroadcastTransaction(txn *Transaction) {
	for peer, status := range bc.Peers {
		if peer != bc.Address && status {
			log.Println("Broadcasting transaction to peer:", peer, "transaction:", txn.ToJson())
			bc.SendTransactionPeer(peer, txn)
			time.Sleep(constants.PEER_LIST_UPDATE_INTERVAL * time.Second)
		}
	}
}

// FetchBlocks: fetches the last N blocks (defined by FETCH_BLOCK_NUMBER constant) from a given address.
// Makes an HTTP GET request to the fetch-consensus-blocks endpoint, reads the response data,
// and unmarshals it into a BlockchainCore struct. Returns a pointer to the blockchain containing
// the fetched blocks and any errors encountered.
func FetchBlocks(address string) (*BlockchainCore, error) {
	log.Println("Fetching last", constants.FETCH_BLOCK_NUMBER, "blocks")
	outURL := fmt.Sprintf("%s/fetch-consensus-blocks", address)
	resp, err := http.Get(outURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bc BlockchainCore
	err = json.Unmarshal(data, &bc)
	if err != nil {
		return nil, err
	}

	return &bc, nil
}

// verifyBlocks: verifies the integrity of a blockchain by checking block hashes and mining difficulty.
// Takes a slice of Block pointers and returns true if all blocks are valid, false otherwise.
// Validates that the genesis block and all subsequent blocks meet the required mining difficulty,
// and that each block's previous hash matches the actual hash of the previous block.
func verifyBlocks(chain []*Block) bool {
	if chain[0].BlockNumber != 0 && chain[0].Hash()[2:2+constants.MINING_DIFFICULTY] != strings.Repeat("0", constants.MINING_DIFFICULTY) {
		log.Println("Chain verification failed for block", chain[0].BlockNumber, "hash:", chain[0].Hash())
		return false
	}

	for i := 1; i < len(chain); i++ {
		if chain[i-1].Hash() != chain[i].PrevHash {
			log.Println("Prev hash verification failed for block", chain[0].BlockNumber)
			return false
		}

		if chain[i].Hash()[2:2+constants.MINING_DIFFICULTY] != strings.Repeat("0", constants.MINING_DIFFICULTY) {
			log.Println("Chain verification failed for block", chain[0].BlockNumber, "hash:", chain[0].Hash())
			return false
		}
	}

	return true
}

// UpdateBlockchain: updates the blockchain with a new chain of blocks. Takes a slice of new blocks,
// updates the blockchain's blocks array by appending the new chain at the correct position based on
// block numbers, and updates the transaction pool by removing any transactions that are now included
// in the blockchain. Thread-safe using mutex locks. After updating, saves the new blockchain state to  the database
func (bc *BlockchainCore) UpdateBlockchain(newChain []*Block) {
	mutex.Lock()
	defer mutex.Unlock()

	blocks := []*Block{}
	initIndex := newChain[0].BlockNumber
	blocks = append(blocks, bc.Blocks[:initIndex]...)
	blocks = append(blocks, newChain...)

	bc.Blocks = blocks

	// Update transaction pool
	found := map[string]bool{}
	for _, txn := range bc.TransactionPool {
		found[txn.TransactionHash] = false
	}

	for _, block := range newChain {
		for _, txn := range block.Transactions {
			_, ok := found[txn.TransactionHash]
			if ok {
				found[txn.TransactionHash] = true
			}
		}

		newTxnPool := []*Transaction{}
		for _, txn := range bc.TransactionPool {
			if !found[txn.TransactionHash] {
				newTxnPool = append(newTxnPool, txn)
			}
		}

		bc.TransactionPool = newTxnPool

		// Save the blockchain to the database
		err := DBAddBlockchain(*bc)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (bc *BlockchainCore) RunConsensus() {
	for {
		log.Println("Running consensus...")
		longestChain := bc.Blocks
		lengthOfTheLongestChain := bc.Blocks[len(bc.Blocks)-1].BlockNumber + 1
		longestChainIsOur := true
		for peer, status := range bc.Peers {
			if peer != bc.Address && status {
				bc1, err := FetchBlocks(peer)
				if err != nil {
					log.Println("Error while  fetching blocks from peer:", peer, "Error:", err.Error())
					continue
				}

				lengthOfTheFetchedChain := bc1.Blocks[len(bc1.Blocks)-1].BlockNumber + 1
				if lengthOfTheFetchedChain > lengthOfTheLongestChain {
					longestChain = bc1.Blocks
					lengthOfTheLongestChain = lengthOfTheFetchedChain
					longestChainIsOur = false
				}
			}
		}

		if longestChainIsOur {
			log.Println("Our chain is the longest, not updating.")
			time.Sleep(constants.CONSENSUS_PAUSE_INTERVAL * time.Second)
			continue
		}

		if verifyBlocks(longestChain) {
			// Stop mining
			bc.MiningLocked = true

			bc.UpdateBlockchain(longestChain)

			// Restart mining
			bc.MiningLocked = false

			log.Println("Blockchain update complete!")
		} else {
			log.Println("Chain Verification Failed, not updating my blockchain")
		}

		time.Sleep(constants.CONSENSUS_PAUSE_INTERVAL * time.Second)
	}

}
