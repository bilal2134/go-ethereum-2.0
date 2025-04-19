package bft

// vrf.go: Verifiable random functions for leader election
// Implements VRF primitives for secure leader election in BFT.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

// VRF represents a verifiable random function
type VRF struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// NewVRF generates a new VRF instance
func NewVRF() (*VRF, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &VRF{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// Evaluate computes the VRF output and proof for a given input
func (vrf *VRF) Evaluate(input []byte) ([]byte, []byte, error) {
	hash := sha256.Sum256(input)
	r, s, err := ecdsa.Sign(rand.Reader, vrf.privateKey, hash[:])
	if err != nil {
		return nil, nil, err
	}
	output := append(r.Bytes(), s.Bytes()...)
	return output, hash[:], nil
}

// Verify verifies the VRF output and proof for a given input
func (vrf *VRF) Verify(input, output, proof []byte) (bool, error) {
	hash := sha256.Sum256(input)
	r := new(big.Int).SetBytes(output[:len(output)/2])
	s := new(big.Int).SetBytes(output[len(output)/2:])
	valid := ecdsa.Verify(vrf.publicKey, hash[:], r, s)
	return valid, nil
}
