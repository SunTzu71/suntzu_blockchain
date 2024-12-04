package blockchainserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
)

type BlockchainServer struct {
	Port          uint64                     `json:"port"`
	BlockchainPtr *blockchain.BlockchainCore `json:"blockchain"`
}

// GetBlockchain: handles HTTP requests to retrieve the blockchain data
// Returns the blockchain as JSON for GET requests and an error for other methods
func (bcs *BlockchainServer) GetBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		io.WriteString(w, bcs.BlockchainPtr.ToJson())
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
}

// GetBalance: handles HTTP requests to retrieve the balance for a given address
// Returns the balance as JSON for GET requests and an error for other methods
func (bcs *BlockchainServer) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		address := r.URL.Query().Get("address")
		x := struct {
			Balance uint64 `json:"balance"`
		}{
			bcs.BlockchainPtr.CalculateTotalCrypto(address),
		}
		mBalance, err := json.Marshal(x)
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, string(mBalance))
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
}

// GetNonRewardedTransactions: handles HTTP requests to retrieve all non-rewarded transactions
// Returns the list of non-rewarded transactions as JSON for GET requests and an error for other methods
func (bcs *BlockchainServer) GetNonRewardedTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		transactionList := bcs.BlockchainPtr.GetAllNonRewardedTransactions()
		bs, err := json.Marshal(transactionList)
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, string(bs))
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
}

// CreateBlockchainServer: creates a new blockchain server with the given port and blockchain reference
func CreateBlockchainServer(port uint64, blockchainPtr *blockchain.BlockchainCore) *BlockchainServer {
	bcs := new(BlockchainServer)
	bcs.Port = port
	bcs.BlockchainPtr = blockchainPtr

	return bcs
}

// StartBlockchainServer: starts the server to handle blockchain requests
func (bcs *BlockchainServer) StartBlockchainServer() {
	http.HandleFunc("/", bcs.GetBlockchain)
	http.HandleFunc("/balance", bcs.GetBalance)
	http.HandleFunc("/get-non-rewarded-transactions", bcs.GetNonRewardedTransactions)

	log.Println("Starting server on port " + strconv.Itoa(int(bcs.Port)))

	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(bcs.Port)), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
