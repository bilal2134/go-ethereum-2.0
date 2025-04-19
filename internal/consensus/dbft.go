package consensus

// dbft.go: Delegated Byzantine Fault Tolerance (dBFT) logic
// Implements dBFT consensus logic for hybrid protocol.

import (
	"fmt"
	"sync"
)

// dBFTState represents the state of a dBFT consensus round.
type dBFTState struct {
	Proposer  string
	Proposal  string
	Votes     map[string]string // nodeID -> vote ("yes"/"no")
	Committed bool
	mu        sync.Mutex
}

// dBFTConsensus manages the dBFT protocol.
type dBFTConsensus struct {
	Nodes        []string
	AuthManager  *AuthManager
	CurrentRound *dBFTState
}

// NewdBFTConsensus creates a new dBFTConsensus instance.
func NewdBFTConsensus(nodes []string, auth *AuthManager) *dBFTConsensus {
	return &dBFTConsensus{
		Nodes:        nodes,
		AuthManager:  auth,
		CurrentRound: nil,
	}
}

// Propose allows an authenticated node to propose a block.
func (d *dBFTConsensus) Propose(proposer, proposal string) error {
	if !d.isAuthenticated(proposer) {
		return fmt.Errorf("node %s not authenticated", proposer)
	}
	d.CurrentRound = &dBFTState{
		Proposer: proposer,
		Proposal: proposal,
		Votes:    make(map[string]string),
	}
	return nil
}

// Vote allows authenticated nodes to vote on the proposal.
func (d *dBFTConsensus) Vote(nodeID, vote string) error {
	if d.CurrentRound == nil {
		return fmt.Errorf("no active round")
	}
	if !d.isAuthenticated(nodeID) {
		return fmt.Errorf("node %s not authenticated", nodeID)
	}
	d.CurrentRound.mu.Lock()
	defer d.CurrentRound.mu.Unlock()
	d.CurrentRound.Votes[nodeID] = vote
	return nil
}

// Commit checks if consensus is reached and commits the proposal if so.
func (d *dBFTConsensus) Commit() (bool, string) {
	if d.CurrentRound == nil {
		return false, ""
	}
	d.CurrentRound.mu.Lock()
	defer d.CurrentRound.mu.Unlock()
	yesVotes := 0
	total := 0
	for _, vote := range d.CurrentRound.Votes {
		total++
		if vote == "yes" {
			yesVotes++
		}
	}
	threshold := len(d.Nodes)*2/3 + 1
	if yesVotes >= threshold {
		d.CurrentRound.Committed = true
		return true, d.CurrentRound.Proposal
	}
	return false, ""
}

// isAuthenticated checks if a node is authenticated and trusted.
func (d *dBFTConsensus) isAuthenticated(nodeID string) bool {
	node, exists := d.AuthManager.nodes[nodeID]
	return exists && node.TrustScore >= 0
}
