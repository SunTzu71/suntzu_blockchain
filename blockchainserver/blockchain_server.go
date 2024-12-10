package blockchainserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
	"github.com/SunTzu71/suntzu_blockchain/constants"
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
		return
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
		return
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
		return
	}
}

// SendTranactionBlockchain: handles HTTP requests to add a new transaction to the blockchain
// Accepts a transaction as JSON in POST requests and returns the added transaction
// Returns an error for other methods
func (bcs *BlockchainServer) SendTranactionBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		defer r.Body.Close()

		request, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()

		var newTransaction blockchain.Transaction
		err = json.Unmarshal(request, &newTransaction)
		if err != nil {
			log.Fatal(err)
		}
		go bcs.BlockchainPtr.AddTransactionToTransactionPool(&newTransaction)
		io.WriteString(w, newTransaction.ToJson())
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}
}

// CreateBlockchainServer: creates a new blockchain server with the given port and blockchain reference
func CreateBlockchainServer(port uint64, blockchainPtr *blockchain.BlockchainCore) *BlockchainServer {
	bcs := new(BlockchainServer)
	bcs.Port = port
	bcs.BlockchainPtr = blockchainPtr

	return bcs
}

func CheckServerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, constants.BLOCKCHAIN_STATUS)
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}
}

// SendPeersList: handles HTTP requests to update the list of blockchain peers
// Accepts a peer list as JSON in POST requests and returns a success message
// Returns an error for other methods or invalid data
func (bcs *BlockchainServer) SendPeersList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		peersMap, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading peers list:", err)
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}

		var peersList map[string]bool
		err = json.Unmarshal(peersMap, &peersList)
		if err != nil {
			log.Println("Error unmarshalling peers list:", err)
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}
		go bcs.BlockchainPtr.UpdatePeers(peersList)
		res := map[string]string{}
		res["success"] = "success"
		x, err := json.Marshal(res)
		if err != nil {
			log.Println("Error marshalling response:", err)
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}
		io.WriteString(w, string(x))
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}
}

// FetchConsensusBlocks: handles HTTP requests to fetch recent blocks for consensus
// Returns the most recent blocks (up to FETCH_BLOCK_NUMBER) as JSON for GET requests
// If fewer blocks exist than FETCH_BLOCK_NUMBER, returns all blocks
// Returns an error for non-GET methods
func (bcs *BlockchainServer) FetchConsensusBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		blocks := bcs.BlockchainPtr.Blocks
		blockchain1 := new(blockchain.BlockchainCore)
		if len(blocks) < constants.FETCH_BLOCK_NUMBER {
			blockchain1.Blocks = blocks
		} else {
			blockchain1.Blocks = blocks[len(blocks)-constants.FETCH_BLOCK_NUMBER:]
		}
		io.WriteString(w, blockchain1.ToJson())
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}
}

// StartBlockchainServer: starts the server to handle blockchain requests
func (bcs *BlockchainServer) StartBlockchainServer() {
	http.HandleFunc("/", bcs.GetBlockchain)
	http.HandleFunc("/balance", bcs.GetBalance)
	http.HandleFunc("/get-non-rewarded-transactions", bcs.GetNonRewardedTransactions)
	http.HandleFunc("/send-transaction", bcs.SendTranactionBlockchain)
	http.HandleFunc("/send-peers-list", bcs.SendPeersList)
	http.HandleFunc("/check-server-status", CheckServerStatus)
	http.HandleFunc("/fetch-consensus-blocks", bcs.FetchConsensusBlocks)

	log.Println("Starting server on port " + strconv.Itoa(int(bcs.Port)))

	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(bcs.Port)), nil) // TODO: place address in config file
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
