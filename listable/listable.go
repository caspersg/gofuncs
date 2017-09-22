package listable

import (
	"sync"
)

// T is *the* abstract type. go doesn't have generics, so the user code needs to cast it
type T interface{}

// Foldl if we define left fold once for a generic list, then a lot of other functions can be derived from it
func Foldl(list []T, init T, f func(result, next T) T) T {
	result := init
	for _, x := range list {
		result = f(result, x)
	}
	return result
}

// Map applies a function to each item inside the list
func Map(list []T, mapFunc func(T) T) []T {
	result := Foldl(list, []T{}, func(result, next T) T {
		return append(result.([]T), mapFunc(next))
	})
	return result.([]T)
}

// Filter returns all the items which pass the filter func
func Filter(list []T, filterFunc func(T) bool) []T {
	result := Foldl(list, []T{}, func(result, next T) T {
		if filterFunc(next) {
			return append(result.([]T), next)
		}
		return result
	})
	return result.([]T)
}

// Length returns the number of items contained in a foldable
func Length(list []T) int {
	return Foldl(list, 0, func(result, next T) T {
		return result.(int) + 1
	}).(int)
}

// All returns true if all items pass the filterFunc
func All(list []T, filterFunc func(T) bool) bool {
	return Foldl(list, true, func(result, next T) T {
		return result.(bool) && filterFunc(next)
	}).(bool)
}

// Any returns true if any of the items pass the filterFunc
func Any(list []T, filterFunc func(T) bool) bool {
	return Foldl(list, false, func(result, next T) T {
		return result.(bool) || filterFunc(next)
	}).(bool)
}

// Concat concatenates the parameters from both lists, in order
func Concat(a, b []T) []T {
	return Foldl(b, a, func(result, next T) T {
		return append(result.([]T), next)
	}).([]T)
}

// // an internal type to store temporary result values
// // these would need to be defined for each use case
// type intAndFoldable struct {
// 	Int int
// 	Foldable
// }

// // Take will return the first n Items in a Foldable
// func Take(list []T, number int) []T {
// 	init := intAndFoldable{Int: 0, Foldable: []T{}}
// 	// result := Foldl(list, init, func(result, next T) T {
// 	result := Foldl(list, init, func(result, next T) T {
// 		// count := result.(intAndFoldable).Int
// 		// previous := result.(intAndFoldable).Foldable
// 		// if count < number {
// 		// 	return intAndFoldable{Int: count + 1, Foldable: previous.Append(next)}
// 		// }
// 		// return result
// 	}).(intAndFoldable)
// 	return result.Foldable
// }

// // Drop will return only the items after the first n Items in a Foldable
// func Drop(list []T, number int) []T {
// 	init := intAndFoldable{Int: 0, Foldable: foldable.Init()}
// 	result := Foldl(list, (init, func(result, next T) T {
// 		count := result.(intAndFoldable).Int
// 		previous := result.(intAndFoldable).Foldable
// 		if count >= number {
// 			return intAndFoldable{Int: count + 1, Foldable: previous.Append(next)}
// 		}
// 		return intAndFoldable{Int: count + 1, Foldable: previous}
// 	}).(intAndFoldable)
// 	return result.Foldable
// }

func newPromise(waitGroup *sync.WaitGroup, mapFunc func() T) *T {
	var p T
	waitGroup.Add(1)
	go func() {
		p = mapFunc()
		waitGroup.Done()
	}()
	return &p
}

// ParMap applies a function in parallel to each item inside the foldable
func ParMap(list []T, mapFunc func(T) T) []T {
	waitGroup := &sync.WaitGroup{}
	init := []*T{}
	pendingResults := Foldl(list, init, func(result, next T) T {
		promise := newPromise(waitGroup, func() T { return mapFunc(next) })
		return append(result.([]*T), promise)
	}).([]*T)
	waitGroup.Wait()

	result := []T{}
	for _, item := range pendingResults {
		result = append(result, *item)
	}
	return result
}
