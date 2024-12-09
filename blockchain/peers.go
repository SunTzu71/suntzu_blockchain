package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
