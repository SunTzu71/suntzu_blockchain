package walletserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/SunTzu71/suntzu_blockchain/wallet"
)

type WalletServer struct {
	Port                  uint16 `json:"port"`
	BlockchainNodeAddress string `json:"blockchain_node_address"`
}

// CreateWalletServer creates a new WalletServer with the given port and blockchain node address
func CreateWalletServer(port uint16, blockchainNodeAddress string) *WalletServer {
	ws := new(WalletServer)
	ws.Port = port
	ws.BlockchainNodeAddress = blockchainNodeAddress
	return ws
}

// CreateNewWallet: handles GET requests to create a new wallet and returns wallet details as JSON
func (ws *WalletServer) CreateNewWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		walletNew, err := wallet.NewWallet()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		x := struct {
			PrivateKeyHex string `json:"private_key_hex"`
			PublicKeyHex  string `json:"public_key_hex"`
			Address       string `json:"address"`
		}{
			PrivateKeyHex: walletNew.GetPrivateKeyHex(),
			PublicKeyHex:  walletNew.GetPublicKeyHex(),
			Address:       walletNew.GetAddress(),
		}
		wbs, err := json.Marshal(x)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		io.WriteString(w, string(wbs))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetTotalCryptoFromWallet: handles GET requests to retrieve the total crypto balance for a wallet address
func (ws *WalletServer) GetTotalCryptoFromWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		params := url.Values{}
		params.Add("address", r.URL.Query().Get("address"))
		ourUrl := fmt.Sprintf("%s?%s", ws.BlockchainNodeAddress+"/balance", params.Encode())
		response, err := http.Get(ourUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer response.Body.Close()
		data, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		io.WriteString(w, string(data))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StartWalletServer: initializes and starts the wallet server, setting up HTTP handlers and listening for connections
func (ws *WalletServer) StartWalletServer() {
	http.HandleFunc("/total-from-wallet", ws.GetTotalCryptoFromWallet)
	http.HandleFunc("/create-new-wallet", ws.CreateNewWallet)

	log.Printf("Wallet server listening on port %d", ws.Port)

	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(ws.Port)), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
