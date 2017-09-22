package foldable

import (
	"sort"
)

// Hash / Map / Dictionary
// key is always a key, so we can sort them
type Hash struct {
	Value map[string]T
}

// HashEntry is the key value pair
type HashEntry struct {
	Key   string
	Value T
}

func (foldable Hash) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	for _, key := range foldable.sortedKeys() {
		result = foldFunc(result, HashEntry{Key: key, Value: foldable.Value[key]})
	}
	return result
}

func (foldable Hash) Init() Foldable {
	return Hash{Value: make(map[string]T)}
}

func (foldable Hash) Append(item T) Foldable {
	foldable.Value[item.(HashEntry).Key] = item.(HashEntry).Value
	return Hash{Value: foldable.Value}
}

func (foldable Hash) sortedKeys() []string {
	// take and drop depend on order, so we need a guaranteed order
	var keys []string
	for k := range foldable.Value {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
