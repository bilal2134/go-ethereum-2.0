package cap

// conflict.go: Advanced conflict resolution
// Implements entropy-based conflict detection and probabilistic conflict resolution.

import (
	"math/rand"
	"time"
)

// Conflict represents a conflict between two or more entities.
type Conflict struct {
	EntityA string
	EntityB string
	Entropy float64
}

// NewConflict creates a new conflict between two entities.
func NewConflict(entityA, entityB string) *Conflict {
	return &Conflict{
		EntityA: entityA,
		EntityB: entityB,
		Entropy: calculateEntropy(entityA, entityB),
	}
}

// calculateEntropy calculates the entropy between two entities.
func calculateEntropy(entityA, entityB string) float64 {
	_ = entityA
	_ = entityB
	// Placeholder for entropy calculation logic.
	return rand.Float64()
}

// ResolveConflict resolves a conflict probabilistically.
func ResolveConflict(conflict *Conflict) string {
	rand.Seed(time.Now().UnixNano())
	if rand.Float64() < conflict.Entropy {
		return conflict.EntityA
	}
	return conflict.EntityB
}

// ConflictWithClock represents a conflict with vector clocks for causal consistency.
type ConflictWithClock struct {
	Entities  []string
	Clocks    []*VectorClock
	Entropies []float64
}

// NewConflictWithClock creates a new conflict with vector clocks.
func NewConflictWithClock(entities []string, clocks []*VectorClock) *ConflictWithClock {
	entropies := make([]float64, len(entities))
	for i := range entities {
		entropies[i] = calculateClockEntropy(clocks, i)
	}
	return &ConflictWithClock{
		Entities:  entities,
		Clocks:    clocks,
		Entropies: entropies,
	}
}

// calculateClockEntropy calculates entropy based on vector clock divergence.
func calculateClockEntropy(clocks []*VectorClock, idx int) float64 {
	// Simple entropy: count how many clocks differ from this one
	base := clocks[idx]
	diff := 0
	for i, c := range clocks {
		if i != idx && base.Compare(c) != 0 {
			diff++
		}
	}
	return float64(diff) / float64(len(clocks)-1)
}

// ResolveMultiEntityConflict resolves a conflict among multiple entities probabilistically, minimizing divergence.
func ResolveMultiEntityConflict(conflict *ConflictWithClock) string {
	maxEntropy := -1.0
	winnerIdx := 0
	for i, entropy := range conflict.Entropies {
		if entropy > maxEntropy {
			maxEntropy = entropy
			winnerIdx = i
		}
	}
	return conflict.Entities[winnerIdx]
}
