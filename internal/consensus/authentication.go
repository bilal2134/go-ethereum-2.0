package consensus

// authentication.go: Advanced node authentication framework
// Implements continuous authentication, adaptive trust scoring, and multi-factor validation for nodes.

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/bilal2134/Blockchain_A3/internal/bft"
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
	nodes      map[string]*Node
	challenges map[string]string
}

// NewAuthManager creates a new AuthManager.
func NewAuthManager() *AuthManager {
	return &AuthManager{
		nodes:      make(map[string]*Node),
		challenges: make(map[string]string),
	}
}

// AddNode adds a new node to the AuthManager.
func (am *AuthManager) AddNode(id, publicKey string) {
	am.nodes[id] = &Node{
		ID:           id,
		PublicKey:    publicKey,
		TrustScore:   0.0,
		LastAuthTime: time.Now(),
	}
	// Initialize challenge map entry
	am.challenges[id] = ""
}

// RemoveNode removes a node from the AuthManager.
func (am *AuthManager) RemoveNode(id string) {
	delete(am.nodes, id)
	delete(am.challenges, id)
}

// NewChallenge generates a random nonce challenge for the node to sign.
func (am *AuthManager) NewChallenge(id string) (string, error) {
	_, exists := am.nodes[id]
	if !exists {
		return "", errors.New("node not found")
	}
	// Generate random 16-byte nonce
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}
	challenge := hex.EncodeToString(nonce)
	am.challenges[id] = challenge
	return challenge, nil
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

// ValidateResponse verifies the HMAC-signed challenge and updates trust score.
func (am *AuthManager) ValidateResponse(id, response string) error {
	node, exists := am.nodes[id]
	if !exists {
		return errors.New("node not found")
	}
	challenge, ok := am.challenges[id]
	if !ok || challenge == "" {
		return errors.New("no active challenge")
	}
	// Verify HMAC using node's PublicKey as key
	if !bft.VerifyHMAC([]byte(node.PublicKey), []byte(challenge), response) {
		// Reduce trust on failure
		node.TrustScore -= 1.0
		return errors.New("HMAC validation failed")
	}
	// Increase trust on success and record time
	node.TrustScore += 1.0
	node.LastAuthTime = time.Now()
	// Clear challenge
	am.challenges[id] = ""
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

// GetTrustScores returns a map of node IDs to their trust scores.
func (am *AuthManager) GetTrustScores() map[string]float64 {
	scores := make(map[string]float64, len(am.nodes))
	for id, node := range am.nodes {
		scores[id] = node.TrustScore
	}
	return scores
}
