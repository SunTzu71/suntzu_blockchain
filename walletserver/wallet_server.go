package walletserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/SunTzu71/suntzu_blockchain/wallet"
)

type WalletServer struct {
	Port                  uint16 `json:"port"`
	BlockchainNodeAddress string `json:"blockchain_node_address"`
}

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

func (ws *WalletServer) StartWalletServer() {
	http.HandleFunc("/create-new-wallet", ws.CreateNewWallet)

	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(ws.Port)), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
