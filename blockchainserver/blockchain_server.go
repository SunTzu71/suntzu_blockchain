package blockchainserver

import (
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

// CreateBlockchainServer: creates a new blockchain server with the given port and blockchain reference
func CreateBlockchainServer(port uint64, blockchainPtr *blockchain.BlockchainCore) *BlockchainServer {
	bcs := new(BlockchainServer)
	bcs.Port = port
	bcs.BlockchainPtr = blockchainPtr

	return bcs
}

// StartBlockchainServer: starts the server to handle blockchain requests
func (bsc *BlockchainServer) StartBlockchainServer() {
	http.HandleFunc("/", bsc.GetBlockchain)
	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(bsc.Port)), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// GetBlockchain: handles HTTP requests to retrieve the blockchain data
// Returns the blockchain as JSON for GET requests and an error for other methods
func (bsc *BlockchainServer) GetBlockchain(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, bsc.BlockchainPtr.ToJson())
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
}
