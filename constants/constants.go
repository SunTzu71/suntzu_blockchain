package constants

// Database path for the blockchain
var BLOCKCHAIN_DB_PATH string

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
	BLOCKCHAIN_KEY             = "blockchain_key"
	ADDRESS_PREFIX             = "suntzuchain"
	TRANSACTION_VERIFY_SUCCESS = "verification_success"
	TRANSACTION_VERIFY_FAILED  = "verification_failed"
	BLOCKCHAIN_STATUS          = "running"
	PEER_LIST_UPDATE_INTERVAL  = 1  // in seconds
	PEER_PING_INTERVAL         = 60 // in seconds
	FETCH_BLOCK_NUMBER         = 50 // number of blocks to fetchfor consensus
	CONSENSUS_PAUSE_INTERVAL   = 10 // in seconds
)
