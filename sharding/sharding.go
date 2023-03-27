package sharding

import (
	"crypto/sha1"
	"sync"
)

// Shard represents a part of a larger data structure
// which has been split to optimize concurrent acces
type Shard[V any] struct {
	sync.RWMutex
	m map[string]V
}

// ShardedMap represents the data structure itself
type ShardedMap[V any] []*Shard[V]

// NewShardedMap creates a new ShardedMap consisting of nshards Shard, allowing to
// concurrently update and/or read from each one of them
func NewShardedMap[V any](nshards uint) ShardedMap[V] {
	shards := make([]*Shard[V], nshards)

	var i uint
	for ; i < nshards; i++ {
		shard := make(map[string]V)
		shards[i] = &Shard[V]{m: shard}
	}

	return shards
}

// getShardIndex computes the index which will be used to retrieve the Shard
// responsible for a given key.
//
// This implementations relies on using a single byte to determine the index,
// therefore only up to index 255 (2^8) shards are available
func (m ShardedMap[V]) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[17]) // arbitrary byte
	return hash % len(m)
}

// getShard returns a pointer to the Shard responsible for the given key
func (m ShardedMap[V]) getShard(key string) *Shard[V] {
	index := m.getShardIndex(key)
	return m[index]
}

// Get returns the value associated with a given key
func (m ShardedMap[V]) Get(key string) V {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

// Set inserts a key with the given value, or updates it
func (m ShardedMap[V]) Set(key string, value V) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

// Contains return a boolean representing whether the given key
// is present in the ShardedMap or not
func (m ShardedMap[V]) Contains(key string) bool {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	_, ok := shard.m[key]
	return ok
}

// Delete removes the given key from the ShardedMap,
// or results in a no-op if it was not present
func (m ShardedMap[V]) Delete(key string) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.m, key)
}

// Keys returns a slice of strings representing all the keys in a ShardedMap.
func (m ShardedMap[V]) Keys() []string {
	keys := make([]string, 0)
	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard[V]) {
			s.RLock()

			for key := range s.m {
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}

			s.RUnlock()
			wg.Done()
		}(shard)
	}

	wg.Wait()
	return keys
}
