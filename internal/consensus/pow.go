package consensus

// pow.go: Proof of Work randomness injection
// Provides PoW-based randomness for hybrid consensus.

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomness generates a random value using PoW.
func GenerateRandomness() (*big.Int, error) {
	max := new(big.Int).Lsh(big.NewInt(1), 256) // 2^256
	randomValue, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	return randomValue, nil
}
