package amf

// Proof.go: Merkle proof generation and compression

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// Node represents a node in the Adaptive Merkle Forest and Merkle trees.
// (Moved to forest.go, use that definition here)

// NewNode creates a new node with the given hash
func NewNode(hash []byte) *Node {
	return &Node{Hash: hash}
}

// MerkleTree represents the Merkle tree
// Use Node for tree nodes
type MerkleTree struct {
	Root *Node
}

// NewMerkleTree creates a new Merkle tree from a list of data
func NewMerkleTree(data [][]byte) (*MerkleTree, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	var nodes []*Node
	for _, datum := range data {
		hash := sha256.Sum256(datum)
		nodes = append(nodes, NewNode(hash[:]))
	}

	for len(nodes) > 1 {
		var newLevel []*Node
		for i := 0; i < len(nodes); i += 2 {
			if i+1 == len(nodes) {
				newLevel = append(newLevel, nodes[i])
			} else {
				left := nodes[i]
				right := nodes[i+1]
				combined := append(left.Hash, right.Hash...)
				combinedHash := sha256.Sum256(combined)
				newLevel = append(newLevel, &Node{
					Left:  left,
					Right: right,
					Hash:  combinedHash[:],
				})
			}
		}
		nodes = newLevel
	}

	return &MerkleTree{Root: nodes[0]}, nil
}

// GenerateProof generates a Merkle proof for a given data item
func (mt *MerkleTree) GenerateProof(data []byte) ([][]byte, error) {
	hash := sha256.Sum256(data)
	hashBytes := hash[:]

	var proof [][]byte
	node := mt.Root
	for node.Left != nil && node.Right != nil {
		if bytes.Equal(node.Left.Hash, hashBytes) || bytes.Equal(node.Right.Hash, hashBytes) {
			if bytes.Equal(node.Left.Hash, hashBytes) {
				proof = append(proof, node.Right.Hash)
				node = node.Left
			} else {
				proof = append(proof, node.Left.Hash)
				node = node.Right
			}
			hashBytes = node.Hash
		} else {
			return nil, errors.New("data not found in the tree")
		}
	}

	return proof, nil
}

// CompressProof compresses a Merkle proof
func CompressProof(proof [][]byte) string {
	if len(proof) == 0 {
		return ""
	}
	return hex.EncodeToString(proof[0])
}

// ProbabilisticProofCompression compresses a Merkle proof probabilistically (stub).
func ProbabilisticProofCompression(proof [][]byte) string {
	// TODO: Implement succinct encoding or hash folding for proof compression
	return CompressProof(proof) // Placeholder: use basic compression for now
}

// AMQFilter is an interface for Approximate Membership Query filters (e.g., Bloom filter).
type AMQFilter interface {
	Add(item string)
	Contains(item string) bool
}

// GenerateAMQProof generates a proof using an AMQ filter (stub).
func GenerateAMQProof(filter AMQFilter, item string) bool {
	return filter.Contains(item)
}

// AccumulatorProof uses a cryptographic accumulator for proof size reduction (stub).
func AccumulatorProof(acc interface{}, item string) bool {
	// TODO: Integrate with accumulator.go for real cryptographic accumulator
	return true // Placeholder
}
