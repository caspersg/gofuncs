package funcs

import (
	"sync"
)

//
// go does not support generics
// code generation is an option, (like https://github.com/cheekybits/genny) but it's an added complication
// the core of the problem is to be able to re-use code
// so what language features does go have to allow code re-use
// interfaces.
// so I wanted to explore some basic functional list processing and see what could be done

// interfaces can only be defined on structs
// so every type will need to be wrapped in a struct
// the implementation will need to be responsible for how the actual value is wrapped/unwrapped
// the call site also needs to wrap and unwrap as well as cast
// this is obviously sacrificing a fair bit of type safety
// but the casts are restricted to functions you pass in and the overall result
// and the user should be aware of what those types are

// T could be a single value, like an int, or another Foldable
// it takes the place of the generic type
type T interface{}

// Foldable this needs to be implemented for each specific type
type Foldable interface {
	// the main way to process the Items within a Foldable
	Foldl(init T, f func(result, next T) T) T
	// there needs to be a way to create an empty version
	Init() Foldable
	// there needs to be a way to combine an Item and a Foldable
	Append(item T) Foldable
}

// there's a few things that can be defined with just a (left) fold
// the interface Foldable *cannot* be the receiver of the function, but that just shows that it can work with any type

// Map applies a function to each item inside the foldable
func Map(foldable Foldable, mapFunc func(T) T) Foldable {
	result := foldable.Foldl(foldable.Init(), func(result, next T) T {
		return result.(Foldable).Append(mapFunc(next))
	}).(Foldable)
	return result
}

// Filter returns all the items which pass the filter func
func Filter(foldable Foldable, filterFunc func(T) bool) Foldable {
	result := foldable.Foldl(foldable.Init(), func(result, next T) T {
		if filterFunc(next) {
			return result.(Foldable).Append(next)
		}
		return result
	}).(Foldable)
	return result
}

// some generic functions operate on int values, so we need to define an internal intItem type.
type intItem struct {
	Value int
}

// Length returns the number of items contained in a foldable
func Length(foldable Foldable) int {
	count := intItem{Value: 0}
	result := foldable.Foldl(count, func(result, next T) T {
		return intItem{Value: result.(intItem).Value + 1}
	}).(intItem)
	return result.Value
}

// some generic functions operate on boolean values, so we need to define an internal boolItem type.
type boolItem struct {
	Value bool
}

// All returns true if all items pass the filterFunc
func All(foldable Foldable, filterFunc func(T) bool) bool {
	result := foldable.Foldl(boolItem{Value: true}, func(result, next T) T {
		return boolItem{Value: result.(boolItem).Value && filterFunc(next)}
	}).(boolItem)
	return result.Value
}

// Any returns true if any of the items pass the filterFunc
func Any(foldable Foldable, filterFunc func(T) bool) bool {
	result := foldable.Foldl(boolItem{Value: false}, func(result, next T) T {
		return boolItem{Value: (result.(boolItem).Value || filterFunc(next))}
	}).(boolItem)
	return result.Value
}

// Concat concatenates the parameters
func Concat(a, b Foldable) Foldable {
	result := b.Foldl(a, func(result, next T) T {
		return result.(Foldable).Append(next)
	}).(Foldable)
	return result
}

// an internal type to store temporary result values
// these would need to be defined for each use case
type intAndFoldable struct {
	Int int
	Foldable
}

// Take will return the first n Items in a Foldable
func Take(foldable Foldable, number int) Foldable {
	init := intAndFoldable{Int: 0, Foldable: foldable.Init()}
	result := foldable.Foldl(init, func(result, next T) T {
		count := result.(intAndFoldable).Int
		previous := result.(intAndFoldable).Foldable
		if count < number {
			return intAndFoldable{Int: count + 1, Foldable: previous.Append(next)}
		}
		return result
	}).(intAndFoldable)
	return result.Foldable
}

// Drop will return only the items after the first n Items in a Foldable
func Drop(foldable Foldable, number int) Foldable {
	init := intAndFoldable{Int: 0, Foldable: foldable.Init()}
	result := foldable.Foldl(init, func(result, next T) T {
		count := result.(intAndFoldable).Int
		previous := result.(intAndFoldable).Foldable
		if count >= number {
			return intAndFoldable{Int: count + 1, Foldable: previous.Append(next)}
		}
		return intAndFoldable{Int: count + 1, Foldable: previous}
	}).(intAndFoldable)
	return result.Foldable
}

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
func ParMap(foldable Foldable, mapFunc func(T) T) Foldable {
	waitGroup := &sync.WaitGroup{}
	init := []*T{}
	pendingResults := foldable.Foldl(init, func(result, next T) T {
		promise := newPromise(waitGroup, func() T { return mapFunc(next) })
		return append(result.([]*T), promise)
	}).([]*T)
	waitGroup.Wait()

	result := foldable.Init()
	for _, item := range pendingResults {
		result = result.Append(*item)
	}
	return result.(Foldable)
}
