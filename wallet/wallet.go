package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/SunTzu71/suntzu_blockchain/blockchain"
	"github.com/SunTzu71/suntzu_blockchain/constants"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey `json:"private_key"`
	PublicKey  *ecdsa.PublicKey  `json:"public_key"`
}

// NewWallet creates and returns a new wallet with a randomly generated private key using elliptic curve P-256.
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	wallet := new(Wallet)
	wallet.PrivateKey = privateKey
	wallet.PublicKey = &privateKey.PublicKey

	return wallet, nil
}

// NewWalletFromPrivateKeyHex creates a new wallet from a hexadecimal private key string.
func NewWalletFromPrivateKeyHex(privateKeyHex string) *Wallet {
	pk := privateKeyHex[2:]
	d := new(big.Int)
	d.SetString(pk, 16)

	var npk ecdsa.PrivateKey
	npk.D = d
	npk.PublicKey.Curve = elliptic.P256()

	// TODO: replace this with crypto/ecdh package
	npk.PublicKey.X, npk.PublicKey.Y = npk.PublicKey.Curve.ScalarBaseMult(d.Bytes())

	wallet := new(Wallet)
	wallet.PrivateKey = &npk
	wallet.PublicKey = &npk.PublicKey

	return wallet
}

// GetPrivateKeyHex returns the private key as a hexadecimal string prefixed with "0x"
func (w *Wallet) GetPrivateKeyHex() string {
	return fmt.Sprintf("0x%x", w.PrivateKey.D)
}

// GetPublicKeyHex returns the public key as a hexadecimal string prefixed with "0x",
// concatenating the X and Y coordinates of the public key point on the curve
func (w *Wallet) GetPublicKeyHex() string {
	return fmt.Sprintf("0x%x%x", w.PublicKey.X, w.PublicKey.Y)
}

// GetAddress generates a unique address for the wallet by:
// 1. Taking the public key (without "0x" prefix)
// 2. Computing its SHA256 hash
// 3. Converting hash to hex string
// 4. Taking last 40 chars and prepending ADDRESS_PREFIX
func (w *Wallet) GetAddress() string {
	hash := sha256.Sum256([]byte(w.GetPublicKeyHex()[2:]))
	hex := fmt.Sprintf("%x", hash[:])
	address := constants.ADDRESS_PREFIX + hex[len(hex)-40:]
	return address
}

// GetSignedTxn takes an unsigned transaction and returns a signed copy of it.
// It does this by:
// 1. Marshaling the transaction to JSON and computing its SHA256 hash
// 2. Signing the hash with the wallet's private key using ECDSA
// 3. Creating a new transaction with the same fields plus signature and public key
// 4. Returns pointer to signed transaction and any error that occurred
func (w *Wallet) GetSignedTransaction(unsignedTxn blockchain.Transaction) (*blockchain.Transaction, error) {
	bs, err := json.Marshal(unsignedTxn)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(bs)

	sig, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	var signedTxn blockchain.Transaction
	signedTxn.From = unsignedTxn.From
	signedTxn.To = unsignedTxn.To
	signedTxn.Data = unsignedTxn.Data
	signedTxn.Status = unsignedTxn.Status
	signedTxn.Value = unsignedTxn.Value
	signedTxn.Timestamp = unsignedTxn.Timestamp
	signedTxn.TransactionHash = unsignedTxn.TransactionHash

	signedTxn.Signature = sig
	signedTxn.PublicKey = w.GetPublicKeyHex()

	return &signedTxn, nil
}
