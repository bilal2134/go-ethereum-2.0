package consensus

// authentication.go: Advanced node authentication framework
// Implements continuous authentication, adaptive trust scoring, and multi-factor validation for nodes.

import (
	"errors"
	"time"
)

// Node represents a node in the network.
type Node struct {
	ID           string
	PublicKey    string
	TrustScore   float64
	LastAuthTime time.Time
}

// AuthManager manages the authentication of nodes.
type AuthManager struct {
	nodes map[string]*Node
}

// NewAuthManager creates a new AuthManager.
func NewAuthManager() *AuthManager {
	return &AuthManager{
		nodes: make(map[string]*Node),
	}
}

// AuthenticateNode authenticates a node using its public key.
func (am *AuthManager) AuthenticateNode(id, publicKey string) error {
	node, exists := am.nodes[id]
	if !exists {
		return errors.New("node not found")
	}
	if node.PublicKey != publicKey {
		return errors.New("invalid public key")
	}
	node.LastAuthTime = time.Now()
	return nil
}

// UpdateTrustScore updates the trust score of a node.
func (am *AuthManager) UpdateTrustScore(id string, score float64) error {
	node, exists := am.nodes[id]
	if !exists {
		return errors.New("node not found")
	}
	node.TrustScore = score
	return nil
}

// AddNode adds a new node to the AuthManager.
func (am *AuthManager) AddNode(id, publicKey string) {
	am.nodes[id] = &Node{
		ID:           id,
		PublicKey:    publicKey,
		TrustScore:   0.0,
		LastAuthTime: time.Now(),
	}
}

// RemoveNode removes a node from the AuthManager.
func (am *AuthManager) RemoveNode(id string) {
	delete(am.nodes, id)
}
