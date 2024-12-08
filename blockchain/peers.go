package blockchain

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	bc.Peers = peers

	err := DBAddBlockchain(*bc)
	if err != nil {
		log.Fatal(err)
	}
}
