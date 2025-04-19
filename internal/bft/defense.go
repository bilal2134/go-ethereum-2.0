package bft

// defense.go: Cryptographic defensive mechanisms against attack vectors
// Implements cryptographic defenses for BFT.

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HMAC generates a keyed-hash message authentication code using SHA-256.
func HMAC(key, message []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC verifies the HMAC of a message.
func VerifyHMAC(key, message []byte, messageMAC string) bool {
	expectedMAC := HMAC(key, message)
	return hmac.Equal([]byte(expectedMAC), []byte(messageMAC))
}

// DefenseManager orchestrates multi-layer BFT defensive checks.
type DefenseManager struct {
	rep *ReputationSystem
	vrf *VRF
}

// NewDefenseManager creates a DefenseManager with reputation system and VRF.
func NewDefenseManager(rep *ReputationSystem, vrf *VRF) *DefenseManager {
	return &DefenseManager{rep: rep, vrf: vrf}
}

// EvaluateState runs layered defenses: ZKP, VRF output verification, and MPC consensus stub.
// Returns true if state is considered valid.
func (dm *DefenseManager) EvaluateState(nodeID string, statement string, zk *ZKProof, randomness []byte, mpcResult *MPCProtocol) bool {
	// 1. Zero-knowledge proof verification
	if !VerifyZKProof(zk) {
		dm.rep.UpdateReputation(nodeID, -10)
		return false
	}
	// 2. VRF verification: verify that randomness matches VRF proof
	ok, err := dm.vrf.Verify([]byte(statement), randomness, zk.ProofData)
	if err != nil || !ok {
		dm.rep.UpdateReputation(nodeID, -5)
		return false
	}
	// 3. MPC consensus stub: ensure mpcResult is non-nil
	if mpcResult == nil || len(mpcResult.Participants) == 0 {
		dm.rep.UpdateReputation(nodeID, -2)
		return false
	}
	// On full success, boost reputation slightly
	dm.rep.UpdateReputation(nodeID, 1)
	return true
}
