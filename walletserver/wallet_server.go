package walletserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
	"github.com/SunTzu71/suntzu_blockchain/constants"
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

// SendTransaction: handles POST requests to create and send a new transaction using the provided private key
// and transaction details, sending it to the blockchain node and returning the response
// TODO: Find a better way to send transaction and not send private key
func (ws *WalletServer) SendTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		privateKey := r.URL.Query().Get("privateKey")
		dataBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer r.Body.Close()

		var trans1 blockchain.Transaction
		err = json.Unmarshal(dataBytes, &trans1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		wallet1 := wallet.NewWalletFromPrivateKeyHex(privateKey)
		myTransaction := blockchain.NewTransaction(wallet1.GetAddress(), trans1.To, trans1.Value, []byte{})
		myTransaction.Status = constants.PENDING
		newTransaction, err := wallet1.GetSignedTransaction(*myTransaction)

		newTransactionBytes, err := json.Marshal(newTransaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Send transaction to blockchain
		response, err := http.Post(ws.BlockchainNodeAddress+"/send-transaction", "application/json", bytes.NewBuffer(newTransactionBytes))
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
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}

// StartWalletServer: initializes and starts the wallet server, setting up HTTP handlers and listening for connections
func (ws *WalletServer) StartWalletServer() {
	http.HandleFunc("/total-from-wallet", ws.GetTotalCryptoFromWallet)
	http.HandleFunc("/create-new-wallet", ws.CreateNewWallet)
	http.HandleFunc("/send-wallet-transaction", ws.SendTransaction)

	log.Printf("Wallet server listening on port %d", ws.Port)

	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(ws.Port)), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
