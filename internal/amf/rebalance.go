package amf

// Rebalance.go: Shard splitting/merging and rebalancing

import (
	"log"
)

// RebalanceConfig holds thresholds for shard rebalancing.
type RebalanceConfig struct {
	SplitThreshold int // Load above which a shard is split
	MergeThreshold int // Combined load below which shards are merged
}

// RebalanceForest checks all shards and splits/merges as needed.
func RebalanceForest(f *Forest, cfg RebalanceConfig) {
	// Split overloaded shards
	for id, shard := range f.Shards {
		if shard.Load > cfg.SplitThreshold {
			log.Printf("Splitting shard %d (load=%d)", id, shard.Load)
			f.SplitShard(id, cfg.SplitThreshold)
		}
	}
	// Merge underutilized shards (naive pairwise, can be improved)
	ids := f.DiscoverShardIDs()
	used := make(map[int]bool)
	for i := 0; i < len(ids); i++ {
		if used[ids[i]] {
			continue
		}
		for j := i + 1; j < len(ids); j++ {
			if used[ids[j]] {
				continue
			}
			shard1, _ := f.GetShard(ids[i])
			shard2, _ := f.GetShard(ids[j])
			if shard1.Load+shard2.Load < cfg.MergeThreshold {
				log.Printf("Merging shards %d and %d (loads=%d,%d)", ids[i], ids[j], shard1.Load, shard2.Load)
				f.MergeShards(ids[i], ids[j], cfg.MergeThreshold)
				used[ids[j]] = true
				break
			}
		}
	}
}

// Optionally, call RebalanceForest after each transaction batch or periodically.
// This maintains adaptive sharding and cryptographic integrity.

// Rebalance is a stub function for shard rebalancing logic.
func Rebalance() {
	// TODO: Implement shard rebalancing logic
}
