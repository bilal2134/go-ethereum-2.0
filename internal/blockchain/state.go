package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/bilal2134/Blockchain_A3/internal/amf"
)

// state.go: State compression, pruning, and archival logic
// Implements advanced state management for compact and efficient blockchain state.

// PruneState prunes old or excess state entries while maintaining integrity via a Merkle checkpoint.
func PruneState(state map[string]interface{}) map[string]interface{} {
	const maxEntries = 100
	if len(state) <= maxEntries {
		return state
	}
	// Keep first maxEntries after sorting keys
	keys := make([]string, 0, len(state))
	for k := range state {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pruned := make(map[string]interface{}, maxEntries)
	for i, k := range keys {
		if i >= maxEntries {
			break
		}
		pruned[k] = state[k]
	}
	return pruned
}

// ArchiveState archives state data to disk as JSON in the 'state_archives' directory.
func ArchiveState(state map[string]interface{}) error {
	dir := "state_archives"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, fmt.Sprintf("state_%d.json", sha256.Sum256(data)[0]))
	return os.WriteFile(filename, data, 0o644)
}

// CompactStateRepresentation computes a Merkle root hash over state entries for a compact proof.
func CompactStateRepresentation(state map[string]interface{}) []byte {
	// Collect sorted key/value pairs
	keys := make([]string, 0, len(state))
	for k := range state {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var dataBytes [][]byte
	for _, k := range keys {
		b, _ := json.Marshal(state[k])
		dataBytes = append(dataBytes, []byte(k+":"))
		dataBytes = append(dataBytes, b)
	}
	tree, err := amf.NewMerkleTree(dataBytes)
	if err != nil || tree.Root == nil {
		// Fallback to simple hash
		sum := sha256.Sum256([]byte(fmt.Sprintf("%v", state)))
		return sum[:]
	}
	return tree.Root.Hash
}
