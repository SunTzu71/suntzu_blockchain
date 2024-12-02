package blockchain

import (
	"encoding/json"

	"github.com/SunTzu71/suntzu_blockchain/constants"
	"github.com/syndtr/goleveldb/leveldb"
)

func DBAddBllockchain(bs BlockchainCore) error {
	db, err := leveldb.OpenFile(constants.BLOCKCHAIN_DB_PATH, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// Save to database
	value, err := json.Marshal(bs)
	if err != nil {
		return err
	}

	err = db.Put([]byte(constants.BLOCKCHAIN_KEY), value, nil)
	if err != nil {
		return err
	}

	return nil
}

func DBGetBlockchain() (*BlockchainCore, error) {
	db, err := leveldb.OpenFile(constants.BLOCKCHAIN_DB_PATH, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	data, err := db.Get([]byte(constants.BLOCKCHAIN_KEY), nil)
	if err != nil {
		return nil, err
	}

	var bs BlockchainCore
	err = json.Unmarshal(data, &bs)
	if err != nil {
		return nil, err
	}

	return &bs, nil
}

func DBKeyExists() bool {
	db, err := leveldb.OpenFile(constants.BLOCKCHAIN_DB_PATH, nil)
	if err != nil {
		return false
	}
	defer db.Close()

	exists, err := db.Has([]byte(constants.BLOCKCHAIN_KEY), nil)
	if err != nil {
		return false
	}

	return exists
}
