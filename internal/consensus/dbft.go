package consensus

// dbft.go: Delegated Byzantine Fault Tolerance (dBFT) logic
// Implements dBFT consensus logic for hybrid protocol.

import (
	"fmt"
	"sync"

	"github.com/bilal2134/Blockchain_A3/internal/bft"
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
	repSystem    *bft.ReputationSystem
}

// NewdBFTConsensus creates a new dBFTConsensus instance with reputation-based thresholding.
func NewdBFTConsensus(nodes []string, auth *AuthManager, rep *bft.ReputationSystem) *dBFTConsensus {
	return &dBFTConsensus{
		Nodes:        nodes,
		AuthManager:  auth,
		CurrentRound: nil,
		repSystem:    rep,
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
	// Dynamic threshold based on node reputations
	// Base threshold is 2/3 of nodes +1, adjusted if average reputation is high
	base := len(d.Nodes)*2/3 + 1
	avgRep := 0
	for _, id := range d.Nodes {
		avgRep += d.repSystem.GetReputation(id)
	}
	if len(d.Nodes) > 0 {
		avgRep = avgRep / len(d.Nodes)
	}
	// If average reputation > 5, tighten threshold to 3/4 of nodes
	threshold := base
	if avgRep > 5 {
		threshold = len(d.Nodes)*3/4 + 1
	}
	if yesVotes >= threshold {
		// Multi-layer defense before commit
		// Prepare proofs
		statement := d.CurrentRound.Proposal
		// Zero-knowledge proof
		zk := bft.GenerateZKProof(statement)
		// VRF for randomness proof
		vrfInst, _ := bft.NewVRF()
		randProof, _, _ := vrfInst.Evaluate([]byte(statement))
		// MPC collective proof
		mpcRes := bft.RunMPC(d.Nodes, []byte(statement))
		// Defense manager
		dm := bft.NewDefenseManager(d.repSystem, vrfInst)
		if !dm.EvaluateState(d.CurrentRound.Proposer, statement, zk, randProof, mpcRes) {
			// Defense failed: abort commit
			return false, ""
		}
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
