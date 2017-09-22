package foldable

import (
	"sync"
)

// I wanted to explore some basic functional list processing and see what could be done
// The idea is to have as much reusable code which can be derived from the smallest amount of unique code that needs to be writtern for each implementation.
// The ideal case is where zero unique code needs to be written for each case, which is why generics is helpful in statically typed languages

// go does not support generics though
// code generation is an option, (like https://github.com/cheekybits/genny) but it's an added complication
// the core of the problem is not being able to re-use code
// so without generics what language features does go have to allow this kind of code re-use
// interfaces.

// interfaces can only be defined on structs
// so every basic value will need to be wrapped in a struct
// the implementation of the interface will need to be responsible for how the actual value is wrapped/unwrapped
// each call site also needs to wrap and unwrap as well as cast the value from the generic type T
// this is obviously sacrificing a fair bit of type safety
// but the casts are restricted to functions you pass in and the overall result
// and the user should be aware of what those types are, as they will more than likely be at the call site

// T could be a single value, like an int, or another Foldable
// T is *the* abstract type.
// go doesn't have generics, so the user code needs to cast to and from it
// this is similar to java before version 1.5, which is when java got generics
// because go also doesn't have covariant type checking, this generic type will also leak a bit further into the code
type T interface{}

// Foldable needs to be implemented for each specific type
// initially I implemented for each container type and value type, like IntListFoldable
// but it can also be implemented for just the container type on a generic value, like List
// the main difference is that the generic list type requires that individual values need to be cast as well
type Foldable interface {
	// the main way to process the Items within a Foldable
	Foldl(init T, f func(result, next T) T) T
	// there needs to be a way to create an empty version
	Init() Foldable
	// there needs to be a way to combine an Item and a Foldable
	Append(item T) Foldable
}

// there's a large number of functions that can be defined with just a (left) fold

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
	Int      int
	Foldable Foldable
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

type foldableTuple struct {
	Pass Foldable
	Fail Foldable
}

// Partition returns a the set of elements which both pass and fail the filter function
func Partition(foldable Foldable, filterFunc func(T) bool) (pass, failed Foldable) {
	result := foldable.Foldl(foldableTuple{foldable.Init(), foldable.Init()}, func(result, next T) T {
		previous := result.(foldableTuple)
		if filterFunc(next) {
			return foldableTuple{Pass: previous.Pass.Append(next), Fail: previous.Fail}
		}
		return foldableTuple{Pass: previous.Pass, Fail: previous.Fail.Append(next)}
	}).(foldableTuple)
	return result.Pass, result.Fail
}
