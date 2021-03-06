package foldable

import (
	"sort"
)

// StringIntEntryItem is the key value pair
type StringIntEntryItem struct {
	Key   string
	Value int
}

type StringIntMapFoldable struct {
	Map map[string]int
}

func (foldable StringIntMapFoldable) sortedKeys() []string {
	// take and drop depend on order, so we need a guaranteed order
	var keys []string
	for k := range foldable.Map {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (foldable StringIntMapFoldable) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	for _, key := range foldable.sortedKeys() {
		result = foldFunc(result, StringIntEntryItem{Key: key, Value: foldable.Map[key]})
	}
	return result
}

func (foldable StringIntMapFoldable) Init() Foldable {
	return StringIntMapFoldable{Map: make(map[string]int)}
}

func (foldable StringIntMapFoldable) Append(item T) Foldable {
	foldable.Map[item.(StringIntEntryItem).Key] = item.(StringIntEntryItem).Value
	return StringIntMapFoldable{Map: foldable.Map}
}
