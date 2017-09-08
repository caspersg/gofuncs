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
// the implementation will need to be responsable for how the actual value is wrapped/unwrapped

// Item could be a single value, like an int, or another Foldable
// it takes the place of the generic type
type Item interface {
}

// FoldableItem is for when an Item also happens to be Foldable
// the generic functions below can return foldables, so we need a conversion method
type FoldableItem interface {
	AsFoldable() Foldable
}

// Foldable this needs to be implemented for each specific type
type Foldable interface {
	// the main way to process the Items within a Foldable
	Foldl(init Item, f func(result, next Item) Item) Item
	// there needs to be a way to create an empty version
	Init() Foldable
	// there needs to be a way to combine an Item and a Foldable
	Append(item Item) Foldable
}

// there's a few things that can be defined with just a (left) fold
// the interface Foldable *cannot* be the receiver of the function, but that just shows that it can work with any type

// Map applies a function to each item inside the foldable
func Map(foldable Foldable, mapFunc func(Item) Item) Foldable {
	result := foldable.Foldl(foldable.Init(), func(result, next Item) Item {
		// TODO just cast to Foldable
		return result.(FoldableItem).AsFoldable().
			Append(mapFunc(next))
	})
	return result.(FoldableItem).AsFoldable()
}

// Filter returns all the items which pass the filter func
func Filter(foldable Foldable, filterFunc func(Item) bool) Foldable {
	result := foldable.Foldl(foldable.Init(), func(result, next Item) Item {
		if filterFunc(next) {
			return result.(FoldableItem).AsFoldable().
				Append(next)
		}
		return result
	})
	return result.(FoldableItem).AsFoldable()
}

// some generic functions operate on int values, so we need to define an internal intItem type.
type intItem struct {
	Value int
}

// Length returns the number of items contained in a foldable
func Length(foldable Foldable) int {
	count := intItem{Value: 0}
	result := foldable.Foldl(count, func(result, next Item) Item {
		return intItem{Value: result.(intItem).Value + 1}
	})
	return result.(intItem).Value
}

// some generic functions operate on boolean values, so we need to define an internal boolItem type.
type boolItem struct {
	Value bool
}

// All returns true if all items pass the filterFunc
func All(foldable Foldable, filterFunc func(Item) bool) bool {
	result := foldable.Foldl(boolItem{Value: true}, func(result, next Item) Item {
		return boolItem{Value: result.(boolItem).Value && filterFunc(next)}
	})
	return result.(boolItem).Value
}

// Any returns true if any of the items pass the filterFunc
func Any(foldable Foldable, filterFunc func(Item) bool) bool {
	result := foldable.Foldl(boolItem{Value: false}, func(result, next Item) Item {
		return boolItem{Value: (result.(boolItem).Value || filterFunc(next))}
	})
	return result.(boolItem).Value
}

// Concat concatenates the parameters
func Concat(a, b Foldable) Foldable {
	result := b.Foldl(a, func(result, next Item) Item {
		return result.(FoldableItem).AsFoldable().Append(next)
	})
	return result.(FoldableItem).AsFoldable()
}

// an internal type to store temporary result values
// these would need to be defined for each use case
type intAndFoldable struct {
	Int int
	Foldable
}

// this method needs to be define, but it doesn't have to be functional as it is not needed
func (i intAndFoldable) AsFoldable() Foldable {
	panic("not supported")
}

// Take will return the first n Items in a Foldable
func Take(foldable Foldable, number int) Foldable {
	init := intAndFoldable{Int: 0, Foldable: foldable.Init()}
	result := foldable.Foldl(init, func(result, next Item) Item {
		count := result.(intAndFoldable).Int
		previous := result.(intAndFoldable).Foldable
		if count < number {
			return intAndFoldable{Int: count + 1, Foldable: previous.Append(next)}
		}
		return result
	})
	return result.(intAndFoldable).Foldable
}

// Drop will return only the items after the first n Items in a Foldable
func Drop(foldable Foldable, number int) Foldable {
	init := intAndFoldable{Int: 0, Foldable: foldable.Init()}
	result := foldable.Foldl(init, func(result, next Item) Item {
		count := result.(intAndFoldable).Int
		previous := result.(intAndFoldable).Foldable
		if count >= number {
			return intAndFoldable{Int: count + 1, Foldable: previous.Append(next)}
		}
		return intAndFoldable{Int: count + 1, Foldable: previous}
	})
	return result.(intAndFoldable).Foldable
}

type resultItem struct {
	Item
}

func NewPromise(waitGroup *sync.WaitGroup, mapFunc func() Item) *resultItem {
	p := &resultItem{}
	waitGroup.Add(1)
	go func() {
		p.Item = mapFunc()
		waitGroup.Done()
	}()
	return p
}

// ParMap applies a function in parallel to each item inside the foldable
func ParMap(foldable Foldable, mapFunc func(Item) Item) Foldable {
	waitGroup := &sync.WaitGroup{}
	init := []*resultItem{}
	pendingResults := foldable.Foldl(init, func(result, next Item) Item {
		promise := NewPromise(waitGroup, func() Item { return mapFunc(next) })
		return append(result.([]*resultItem), promise)
	})
	waitGroup.Wait()

	result := foldable.Init()
	for _, p := range pendingResults.([]*resultItem) {
		result = result.Append(p.Item)
	}
	return result.(FoldableItem).AsFoldable()
}
