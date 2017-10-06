package foldable

import (
	"sort"
)

// Hash / Map / Dictionary
// key is always a key, so we can sort them
type Hash map[string]T

// HashEntry is the key value pair
type HashEntry struct {
	Key   string
	Value T
}

func (foldable Hash) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	for _, key := range foldable.sortedKeys() {
		result = foldFunc(result, HashEntry{Key: key, Value: foldable[key]})
	}
	return result
}

func (foldable Hash) Init() Foldable {
	return make(Hash)
}

func (foldable Hash) Append(item T) Foldable {
	foldable[item.(HashEntry).Key] = item.(HashEntry).Value
	return foldable
}

func (foldable Hash) sortedKeys() []string {
	// take and drop depend on order, so we need a guaranteed order
	var keys []string
	for k := range foldable {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
