package types

// types.go: Common types and interfaces for the blockchain system

// Block represents a single block in the blockchain.
type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

// Blockchain represents the entire chain of blocks.
type Blockchain struct {
	Blocks []Block
}

// Transaction represents a single transaction in the blockchain.
type Transaction struct {
	ID     string
	Amount int
	From   string
	To     string
}

// Wallet represents a user's wallet in the blockchain system.
type Wallet struct {
	Address string
	Balance int
}

// BlockchainInterface defines the methods that a blockchain implementation must have.
type BlockchainInterface interface {
	AddBlock(data string) error
	GetBlock(index int) (Block, error)
	GetBlockchain() []Block
}

// WalletInterface defines the methods that a wallet implementation must have.
type WalletInterface interface {
	CreateWallet() (Wallet, error)
	GetWallet(address string) (Wallet, error)
	Transfer(from string, to string, amount int) error
}
