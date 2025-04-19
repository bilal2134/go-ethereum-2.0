package bft

// reputation.go: Reputation-based node scoring system
// Implements node reputation tracking and scoring for adversarial defense.

type NodeReputation struct {
	NodeID          string
	ReputationScore int
}

type ReputationSystem struct {
	NodeReputations map[string]*NodeReputation
}

func NewReputationSystem() *ReputationSystem {
	return &ReputationSystem{
		NodeReputations: make(map[string]*NodeReputation),
	}
}

func (rs *ReputationSystem) UpdateReputation(nodeID string, scoreDelta int) {
	if nodeRep, exists := rs.NodeReputations[nodeID]; exists {
		nodeRep.ReputationScore += scoreDelta
	} else {
		rs.NodeReputations[nodeID] = &NodeReputation{
			NodeID:          nodeID,
			ReputationScore: scoreDelta,
		}
	}
}

func (rs *ReputationSystem) GetReputation(nodeID string) int {
	if nodeRep, exists := rs.NodeReputations[nodeID]; exists {
		return nodeRep.ReputationScore
	}
	return 0
}
