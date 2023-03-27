package sharding

import "testing"

func TestNewShardedMap(t *testing.T) {
	tests := []struct {
		name    string
		nShards uint
	}{
		{"Empty sharded map", 0},
		{"1 shard", 1},
		{"3 shards", 3},
		{"10 shards", 10},
		{"100 shards", 100},
		{"255 shards", 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMap := NewShardedMap[int](tt.nShards)
			l := len(sMap)

			if l != int(tt.nShards) {
				t.Errorf("Expected len() to return %d - got %d", tt.nShards, l)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	keys := [...]string{"test", "try", "prova", "chance", "again", "123", "QWERTY", "test", "test2", "test3"}

	sMap := NewShardedMap[int](10)
	for i, k := range keys {
		sMap.Set(k, i)
	}

	for i := len(keys) - 1; i >= 0; i-- {
		sMap.Delete(keys[i])
		found := sMap.Contains(keys[i])
		if found {
			t.Errorf("Expected key %s to have been removed, but it was present", keys[i])
		}
	}
}

func TestKeys(t *testing.T) {
	keys := [...]string{"test", "try", "prova", "chance", "again", "123", "QWERTY", "test", "test2", "test3"}

	sMap := NewShardedMap[int](10)
	for i, k := range keys {
		sMap.Set(k, i)
	}

	finalKeys := sMap.Keys()
	for _, k := range keys {
		check := false
		for _, fk := range finalKeys {
			if k == fk {
				check = true
				break
			}
		}
		if !check {
			t.Errorf("Expected Keys() to contain key %s but it was missing", k)
		}
	}
}
