package constants

// Constants used throughout the blockchain
const (
	BLOCKCHAIN_NAME            = "SunTzuChain"
	HEX_PREFIX                 = "0x"
	SUCCESS                    = "success"
	FAILED                     = "failed"
	PENDING                    = "pending"
	MINING_DIFFICULTY          = 5
	MINING_REWARD              = 100 * DECIMAL
	CURRENCY_NAME              = "SZU"
	DECIMAL                    = 100
	BLOCKCHAIN_ADDRESS         = "SunTzu_Faucet"
	BLOCKCHAIN_DB_PATH         = "5001/suntzuchain.db"
	BLOCKCHAIN_KEY             = "blockchain_key"
	ADDRESS_PREFIX             = "suntzuchain"
	TRANSACTION_VERIFY_SUCCESS = "verification_success"
	TRANSACTION_VERIFY_FAILED  = "verification_failed"
	BLOCKCHAIN_STATUS          = "running"
	PEER_LIST_UPDATE_INTERVAL  = 1 // in seconds
)
