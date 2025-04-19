package cap

// vectorclock.go: Causal consistency with vector clocks
// Provides vector clock implementation for causal ordering and conflict detection.

import (
	"sync"
)

// VectorClock represents a vector clock for causal consistency.
type VectorClock struct {
	mu    sync.Mutex
	clock map[string]int
}

// NewVectorClock creates a new vector clock.
func NewVectorClock() *VectorClock {
	return &VectorClock{
		clock: make(map[string]int),
	}
}

// Increment increments the clock for the given node.
func (vc *VectorClock) Increment(node string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.clock[node]++
}

// Update updates the vector clock with another vector clock.
func (vc *VectorClock) Update(other *VectorClock) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	for node, time := range other.clock {
		if time > vc.clock[node] {
			vc.clock[node] = time
		}
	}
}

// Compare compares the vector clock with another vector clock.
// Returns -1 if vc < other, 0 if vc == other, and 1 if vc > other.
func (vc *VectorClock) Compare(other *VectorClock) int {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vcLess := false
	vcGreater := false

	for node, time := range vc.clock {
		otherTime, exists := other.clock[node]
		if !exists {
			vcGreater = true
			continue
		}
		if time < otherTime {
			vcLess = true
		} else if time > otherTime {
			vcGreater = true
		}
		if vcLess && vcGreater {
			return 0
		}
	}

	for node := range other.clock {
		if _, exists := vc.clock[node]; !exists {
			vcLess = true
		}
		if vcLess && vcGreater {
			return 0
		}
	}

	if vcLess {
		return -1
	}
	if vcGreater {
		return 1
	}
	return 0
}
