package bft

// consensus.go: Adaptive consensus thresholds and core BFT consensus logic
// Implements consensus with adaptive thresholds based on node performance.

import (
	"sync"
)

type Node struct {
	ID          string
	Performance int
}

type Consensus struct {
	nodes     []Node
	threshold int
	mu        sync.Mutex
}

func NewConsensus(nodes []Node) *Consensus {
	return &Consensus{
		nodes:     nodes,
		threshold: calculateInitialThreshold(nodes),
	}
}

func calculateInitialThreshold(nodes []Node) int {
	// Implement initial threshold calculation based on node performance
	return len(nodes) / 2
}

func (c *Consensus) UpdateThreshold() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Implement adaptive threshold update logic based on node performance
}

// Integrate with ReputationSystem for adaptive consensus threshold
func (c *Consensus) UpdateThresholdWithReputation(rep *ReputationSystem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Example: Increase threshold if average reputation is high, decrease if low
	totalRep := 0
	for _, node := range c.nodes {
		totalRep += rep.GetReputation(node.ID)
	}
	avgRep := 0
	if len(c.nodes) > 0 {
		avgRep = totalRep / len(c.nodes)
	}
	if avgRep > 10 {
		c.threshold = len(c.nodes)*2/3 + 1 // Stricter consensus
	} else {
		c.threshold = len(c.nodes)/2 + 1 // Default BFT threshold
	}
}

func (c *Consensus) ReachConsensus() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Implement core BFT consensus logic
	return true
}

func (c *Consensus) AddNode(node Node) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nodes = append(c.nodes, node)
	c.UpdateThreshold()
}

func (c *Consensus) RemoveNode(nodeID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, node := range c.nodes {
		if node.ID == nodeID {
			c.nodes = append(c.nodes[:i], c.nodes[i+1:]...)
			break
		}
	}
	c.UpdateThreshold()
}
