package blockchain

// block.go: Advanced block structure for compact state representation and validation
// Includes cryptographic accumulators, multi-level Merkle trees, and entropy-based validation.

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/bilal2134/Blockchain_A3/internal/amf"
)

// Block defines the core blockchain block with compact state and validation fields.
type Block struct {
	Index        int
	Timestamp    time.Time
	PrevHash     string
	Transactions []string
	Accumulator  []byte   // Compact state representation
	MultiMerkle  [][]byte // Multi-level Merkle roots
	Entropy      float64  // Entropy-based validation metric
	Hash         string
}

type Blockchain struct {
	Blocks []*Block
}

// NewBlock creates and initializes a new block with cryptographic accumulator, Merkle root, and entropy.
func NewBlock(index int, prevHash string, txs []string, _ []byte, _ [][]byte, _ float64) *Block {
	b := &Block{
		Index:        index,
		Timestamp:    time.Now(),
		PrevHash:     prevHash,
		Transactions: txs,
	}
	// Accumulator: SHA-256 over concatenated transactions
	data := ""
	for _, tx := range txs {
		data += tx
	}
	accHash := sha256.Sum256([]byte(data))
	b.Accumulator = accHash[:]

	// Level-1 Merkle: transactions
	var dataBytes [][]byte
	for _, tx := range txs {
		dataBytes = append(dataBytes, []byte(tx))
	}
	tree, err := amf.NewMerkleTree(dataBytes)
	if err != nil {
		b.MultiMerkle = [][]byte{}
	} else {
		level1 := tree.Root.Hash
		// Level-2 Merkle: hash(headerHash || level1)
		header := fmt.Sprintf("%d%s%x", index, prevHash, b.Accumulator)
		hdrHash := sha256.Sum256([]byte(header))
		lvl2 := sha256.Sum256(append(hdrHash[:], level1...))
		b.MultiMerkle = [][]byte{level1, lvl2[:]}
		// Derive entropy from hdrHash
		b.Entropy = float64(hdrHash[0]) / 255.0
		// Block hash is second-level root
		b.Hash = fmt.Sprintf("%x", lvl2[:])
	}
	return b
}

// ValidateBlock performs entropy-based and cryptographic validation.
func (b *Block) ValidateBlock() bool {
	// Verify accumulator
	data := ""
	for _, tx := range b.Transactions {
		data += tx
	}
	hash := sha256.Sum256([]byte(data))
	if !bytes.Equal(hash[:], b.Accumulator) {
		return false
	}
	// Recompute level-1 Merkle
	var dataBytes [][]byte
	for _, tx := range b.Transactions {
		dataBytes = append(dataBytes, []byte(tx))
	}
	tree, err := amf.NewMerkleTree(dataBytes)
	if err != nil || len(b.MultiMerkle) < 2 {
		return false
	}
	level1 := tree.Root.Hash
	if !bytes.Equal(level1, b.MultiMerkle[0]) {
		return false
	}
	// Recompute level-2 Merkle
	header := fmt.Sprintf("%d%s%x", b.Index, b.PrevHash, b.Accumulator)
	hdrHash := sha256.Sum256([]byte(header))
	lvl2 := sha256.Sum256(append(hdrHash[:], level1...))
	if !bytes.Equal(lvl2[:], b.MultiMerkle[1]) {
		return false
	}
	// Verify entropy matches header hash
	if b.Entropy != float64(hdrHash[0])/255.0 {
		return false
	}
	// Verify block hash equals second-level root
	if fmt.Sprintf("%x", lvl2[:]) != b.Hash {
		return false
	}
	return true
}

// AddBlock appends a new block to the blockchain.
func (bc *Blockchain) AddBlock(block *Block) {
	bc.Blocks = append(bc.Blocks, block)
}
