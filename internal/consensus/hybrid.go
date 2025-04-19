package consensus

// hybrid.go: Hybrid consensus protocol (PoW randomness + dBFT)
// Implements a novel consensus mechanism combining Proof of Work randomness injection and Delegated BFT principles.

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

// HybridConsensus represents the hybrid consensus protocol.
type HybridConsensus struct {
	mu         sync.Mutex
	round      int
	leader     string
	validators []string
	powRandom  *big.Int
	dbftState  map[string]string
}

// NewHybridConsensus creates a new instance of HybridConsensus.
func NewHybridConsensus(validators []string) *HybridConsensus {
	return &HybridConsensus{
		round:      0,
		validators: validators,
		powRandom:  big.NewInt(0),
		dbftState:  make(map[string]string),
	}
}

// StartRound starts a new round of the hybrid consensus protocol.
func (hc *HybridConsensus) StartRound() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.round++
	hc.powRandom = hc.generatePoWRandomness()
	hc.leader = hc.selectLeader()
	hc.dbftState = make(map[string]string)

	fmt.Printf("Round %d started with leader %s and PoW randomness %s\n", hc.round, hc.leader, hc.powRandom.String())
}

// generatePoWRandomness generates randomness using Proof of Work.
func (hc *HybridConsensus) generatePoWRandomness() *big.Int {
	// Simulate PoW randomness generation
	time.Sleep(1 * time.Second)
	return big.NewInt(time.Now().UnixNano())
}

// selectLeader selects the leader for the current round based on PoW randomness.
func (hc *HybridConsensus) selectLeader() string {
	index := new(big.Int).Mod(hc.powRandom, big.NewInt(int64(len(hc.validators)))).Int64()
	return hc.validators[index]
}

// ProposeBlock allows the leader to propose a new block.
func (hc *HybridConsensus) ProposeBlock(block string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if hc.leader != "" {
		hc.dbftState[hc.leader] = block
		fmt.Printf("Leader %s proposed block: %s\n", hc.leader, block)
	}
}

// Vote allows validators to vote on the proposed block.
func (hc *HybridConsensus) Vote(validator, vote string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if _, exists := hc.dbftState[validator]; exists {
		hc.dbftState[validator] = vote
		fmt.Printf("Validator %s voted: %s\n", validator, vote)
	}
}

// FinalizeRound finalizes the current round and commits the block if consensus is reached.
func (hc *HybridConsensus) FinalizeRound() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Simulate consensus check
	time.Sleep(1 * time.Second)

	// Check if consensus is reached (simple majority for demonstration)
	voteCount := make(map[string]int)
	for _, vote := range hc.dbftState {
		voteCount[vote]++
	}

	for block, count := range voteCount {
		if count > len(hc.validators)/2 {
			fmt.Printf("Consensus reached on block: %s\n", block)
			return
		}
	}

	fmt.Println("Consensus not reached, starting new round.")
	hc.StartRound()
}
