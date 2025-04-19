package amf

// Shard.go: Shard structure and management

type Shard struct {
	ID   int
	Data map[string]interface{}
}

func NewShard(id int) *Shard {
	return &Shard{
		ID:   id,
		Data: make(map[string]interface{}),
	}
}

func (s *Shard) AddData(key string, value interface{}) {
	s.Data[key] = value
}

func (s *Shard) GetData(key string) (interface{}, bool) {
	value, exists := s.Data[key]
	return value, exists
}

func (s *Shard) RemoveData(key string) {
	delete(s.Data, key)
}
