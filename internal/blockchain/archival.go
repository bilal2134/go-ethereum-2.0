package blockchain

// archival.go: Efficient state archival mechanisms
// Provides additional archival utilities for blockchain state.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ArchiveBlock stores a block in an archival storage directory as JSON.
func ArchiveBlock(block interface{}) error {
	blk, ok := block.(*Block)
	if !ok {
		return fmt.Errorf("invalid block type for archival")
	}
	dir := "archives"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(blk, "", "  ")
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, fmt.Sprintf("block_%d.json", blk.Index))
	return os.WriteFile(filename, data, 0o644)
}
