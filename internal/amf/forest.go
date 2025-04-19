package amf

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// Forest.go: Core Adaptive Merkle Forest logic

// Node represents a node in the Adaptive Merkle Forest and Merkle trees.
type Node struct {
	Hash  []byte
	Left  *Node
	Right *Node
}

// ShardWithMeta extends Shard with load tracking and Merkle root.
type ShardWithMeta struct {
	Shard *Shard
	Root  *Node
	Load  int // Computational load (e.g., number of transactions)
	ID    int
	Mutex sync.RWMutex
}

// Forest manages shards and supports dynamic sharding.
type Forest struct {
	Shards map[int]*ShardWithMeta // ShardID -> ShardWithMeta
	Roots  []*Node
	mutex  sync.RWMutex
	// ...other fields as needed...
}

// NewForest creates a new empty Adaptive Merkle Forest.
func NewForest() *Forest {
	return &Forest{
		Shards: make(map[int]*ShardWithMeta),
		Roots:  []*Node{},
	}
}

// AddRoot adds a new root to the forest.
func (f *Forest) AddRoot(root *Node) {
	f.Roots = append(f.Roots, root)
}

// MergeNodes hashes two child nodes into a new parent node for Merkle integrity.
func MergeNodes(left, right *Node) *Node {
	combined := append(left.Hash, right.Hash...)
	h := sha256.Sum256(combined)
	return &Node{Hash: h[:], Left: left, Right: right}
}

// CreateShard creates and adds a new shard to the forest.
func (f *Forest) CreateShard(id int) *ShardWithMeta {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	shard := &ShardWithMeta{
		Shard: NewShard(id),
		ID:    id,
		Load:  0,
	}
	f.Shards[id] = shard
	return shard
}

// GetShard returns a shard by ID (logarithmic time with map).
func (f *Forest) GetShard(id int) (*ShardWithMeta, bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	shard, ok := f.Shards[id]
	return shard, ok
}

// SplitShard splits a shard into two if load exceeds threshold.
func (f *Forest) SplitShard(id int, threshold int) ([]*ShardWithMeta, bool) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	shard, ok := f.Shards[id]
	if !ok || shard.Load < threshold {
		return nil, false
	}
	// Example: split data by even/odd keys (replace with real logic)
	left := &ShardWithMeta{Shard: NewShard(id * 2), ID: id * 2}
	right := &ShardWithMeta{Shard: NewShard(id*2 + 1), ID: id*2 + 1}
	for k, v := range shard.Shard.Data {
		if len(k)%2 == 0 {
			left.Shard.Data[k] = v
		} else {
			right.Shard.Data[k] = v
		}
	}
	// Maintain cryptographic integrity: recompute Merkle roots
	left.Root = BuildMerkleRoot(left.Shard)
	right.Root = BuildMerkleRoot(right.Shard)
	f.Shards[left.ID] = left
	f.Shards[right.ID] = right
	delete(f.Shards, id)
	return []*ShardWithMeta{left, right}, true
}

// MergeShards merges two shards if their combined load is below threshold.
func (f *Forest) MergeShards(id1, id2, threshold int) (*ShardWithMeta, bool) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	shard1, ok1 := f.Shards[id1]
	shard2, ok2 := f.Shards[id2]
	if !ok1 || !ok2 || (shard1.Load+shard2.Load) > threshold {
		return nil, false
	}
	merged := &ShardWithMeta{Shard: NewShard(id1), ID: id1}
	for k, v := range shard1.Shard.Data {
		merged.Shard.Data[k] = v
	}
	for k, v := range shard2.Shard.Data {
		merged.Shard.Data[k] = v
	}
	merged.Load = shard1.Load + shard2.Load
	merged.Root = BuildMerkleRoot(merged.Shard)
	f.Shards[id1] = merged
	delete(f.Shards, id2)
	return merged, true
}

// BuildMerkleRoot computes the Merkle root over sorted shard data key/value pairs.
func BuildMerkleRoot(shard *Shard) *Node {
	// Collect and sort keys for deterministic order
	keys := make([]string, 0, len(shard.Data))
	for k := range shard.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create leaf nodes hashing key and JSON-encoded value
	var nodes []*Node
	for _, k := range keys {
		v := shard.Data[k]
		b, _ := json.Marshal(v)
		leafHash := sha256.Sum256(append([]byte(k+":"), b...))
		nodes = append(nodes, &Node{Hash: leafHash[:]})
	}
	if len(nodes) == 0 {
		// Empty shard: return zero-hash node
		empty := sha256.Sum256(nil)
		return &Node{Hash: empty[:]}
	}
	// Build tree by merging pairs until one root remains
	for len(nodes) > 1 {
		var next []*Node
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				next = append(next, MergeNodes(nodes[i], nodes[i+1]))
			} else {
				// Odd node: carry up
				next = append(next, nodes[i])
			}
		}
		nodes = next
	}
	return nodes[0]
}

// DiscoverShardIDs returns all shard IDs (logarithmic time).
func (f *Forest) DiscoverShardIDs() []int {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	ids := make([]int, 0, len(f.Shards))
	for id := range f.Shards {
		ids = append(ids, id)
	}
	return ids
}

// ReconstructState reconstructs the state from all shards.
func (f *Forest) ReconstructState() map[string]interface{} {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	state := make(map[string]interface{})
	for _, shard := range f.Shards {
		for k, v := range shard.Shard.Data {
			state[k] = v
		}
	}
	return state
}

// AddDataToShard inserts a key/value into the specified shard, updates load and Merkle root, then rebalances the forest.
// Requires a RebalanceConfig to control split/merge thresholds.
func (f *Forest) AddDataToShard(id int, key string, value interface{}, cfg RebalanceConfig) error {
	f.mutex.RLock()
	shardMeta, ok := f.Shards[id]
	f.mutex.RUnlock()
	if !ok {
		return fmt.Errorf("shard %d not found", id)
	}
	shardMeta.Mutex.Lock()
	defer shardMeta.Mutex.Unlock()
	// Add data and update load
	shardMeta.Shard.AddData(key, value)
	shardMeta.Load++
	// Recompute Merkle root
	shardMeta.Root = BuildMerkleRoot(shardMeta.Shard)
	// Trigger rebalance
	RebalanceForest(f, cfg)
	return nil
}
