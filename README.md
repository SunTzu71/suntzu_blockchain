# SunTzu Blockchain

A decentralized blockchain implementation in Go featuring proof-of-work consensus, wallet management, and peer-to-peer networking.

## Features

- Proof-of-work consensus mechanism
- Decentralized peer-to-peer network
- Digital wallet creation and management
- Transaction signing and verification using ECDSA
- Persistent storage using LevelDB
- Mining rewards system
- Real-time blockchain synchronization
- HTTP API for blockchain and wallet interactions

## Architecture

### Core Components

1. **Blockchain Core**
   - Manages blocks, transactions, and mining operations
   - Implements proof-of-work consensus
   - Handles peer-to-peer network communication
   - Validates transactions and maintains the chain state

2. **Wallet System**
   - Creates and manages digital wallets
   - Generates public/private key pairs using ECDSA
   - Signs transactions
   - Manages wallet balances

3. **Network Layer**
   - Handles peer discovery and management
   - Synchronizes blockchain state across nodes
   - Broadcasts transactions and blocks
   - Implements consensus mechanisms

### Data Structures

- **Block**: Contains transactions, timestamps, and chain linking information
- **Transaction**: Records transfer of value between addresses
- **Wallet**: Manages cryptographic keys and addresses
- **BlockchainCore**: Main chain state and operation manager

## Getting Started

### Prerequisites

- Go 1.15 or higher
- LevelDB

### Running the Blockchain

1. Start the first node:
```bash
go run main.go chain -port 8000 -miner <miner_address>
```

2. Start additional nodes:
```bash
go run main.go chain -port 8001 -miner <miner_address> -remote_node http://127.0.0.1:8000
```

3. Start a wallet server:
```bash
go run main.go wallet -port 8080 -node http://127.0.0.1:8000
```

### Using the Launch Script

You can also use the provided launch script to start multiple nodes:
```bash
./launch_blockchain.sh
```

## API Endpoints

### Blockchain Server

- GET `/` - Get full blockchain data
- GET `/balance` - Get address balance
- GET `/get-non-rewarded-transactions` - Get pending transactions
- POST `/send-transaction` - Submit new transaction
- GET `/check-server-status` - Check node status
- GET `/fetch-consensus-blocks` - Get recent blocks for consensus

### Wallet Server

- GET `/create-new-wallet` - Create new wallet
- GET `/total-from-wallet` - Get wallet balance
- POST `/send-wallet-transaction` - Send transaction from wallet

## Consensus Mechanism

The blockchain uses a proof-of-work consensus mechanism where:
- Miners compete to solve computational puzzles
- Mining difficulty is adjusted by required leading zeros
- Successful miners receive rewards in cryptocurrency
- Longest valid chain is accepted as the truth

## Security Features

- ECDSA for transaction signing
- SHA-256 hashing for blocks and transactions
- Balance verification before transaction processing
- Signature verification for all transactions
- Peer verification and validation

## Storage

The blockchain uses LevelDB for persistent storage of:
- Blockchain state
- Blocks
- Transactions
- Peer information

## License

This project is licensed under the MIT License.
