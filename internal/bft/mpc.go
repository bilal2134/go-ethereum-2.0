package bft

import (
	"crypto/sha256"
)

// mpc.go: Multi-party computation protocols for distributed trust
// Implements MPC primitives for distributed trust in BFT.

// MPCProtocol is a stub for a multi-party computation protocol.
type MPCProtocol struct {
	Participants []string
	Result       []byte
}

// RunMPC runs a simple multi-party computation protocol by hashing the inputs and participant IDs.
func RunMPC(participants []string, input []byte) *MPCProtocol {
	// Combine input and participant identifiers
	combined := append([]byte{}, input...)
	for _, id := range participants {
		combined = append(combined, []byte(id)...)
	}
	// Compute a hash as the MPC result
	hash := sha256.Sum256(combined)
	return &MPCProtocol{
		Participants: participants,
		Result:       hash[:],
	}
}
