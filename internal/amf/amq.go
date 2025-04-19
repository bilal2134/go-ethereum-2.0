package amf

// Amq.go: Approximate Membership Query filter logic

// SimpleBloomFilter is a basic Bloom filter implementation for AMQ.
type SimpleBloomFilter struct {
	bits  []bool
	salts []uint32
	size  uint32
}

// NewSimpleBloomFilter creates a new Bloom filter with the given size and salts.
func NewSimpleBloomFilter(size uint32, salts []uint32) *SimpleBloomFilter {
	return &SimpleBloomFilter{
		bits:  make([]bool, size),
		salts: salts,
		size:  size,
	}
}

// hash computes a simple hash for the item with a salt.
func (bf *SimpleBloomFilter) hash(item string, salt uint32) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(item); i++ {
		h ^= uint32(item[i])
		h *= 16777619
	}
	return (h ^ salt) % bf.size
}

// Add inserts an item into the Bloom filter.
func (bf *SimpleBloomFilter) Add(item string) {
	for _, salt := range bf.salts {
		idx := bf.hash(item, salt)
		bf.bits[idx] = true
	}
}

// Contains checks if an item is possibly in the Bloom filter.
func (bf *SimpleBloomFilter) Contains(item string) bool {
	for _, salt := range bf.salts {
		idx := bf.hash(item, salt)
		if !bf.bits[idx] {
			return false
		}
	}
	return true
}
