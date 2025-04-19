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

// ResolveMultiEntityConflict resolves a conflict among multiple entities by choosing the lowest entropy (least divergence).
// If multiple entities tie, uses vector-clock causal ordering to pick the earliest.
func ResolveMultiEntityConflict(conflict *ConflictWithClock) string {
	// Find minimum entropy
	minE := conflict.Entropies[0]
	for _, e := range conflict.Entropies[1:] {
		if e < minE {
			minE = e
		}
	}
	// Collect candidates with min entropy
	candidates := []int{}
	for i, e := range conflict.Entropies {
		if e == minE {
			candidates = append(candidates, i)
		}
	}
	// If only one candidate, choose it
	if len(candidates) == 1 {
		return conflict.Entities[candidates[0]]
	}
	// If two, use vector clock comparison
	// Compare each pair; pick entity whose clock is causally earlier
	best := candidates[0]
	for _, idx := range candidates[1:] {
		cmp := conflict.Clocks[best].Compare(conflict.Clocks[idx])
		if cmp == 1 {
			// current best is later than idx, pick idx
			best = idx
		}
	}
	return conflict.Entities[best]
}

// ResolveConflictWithClock is an alias to ResolveMultiEntityConflict.
var ResolveConflictWithClock = ResolveMultiEntityConflict
