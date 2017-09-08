package funcs

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
	// the generic functions below can return foldables, so we need a conversion method
	AsFoldable() Foldable
}

type Foldable interface {
	// the main way to process the Items within a Foldable
	Foldl(init Item, f func(result, next Item) Item) Item
	// there needs to be a way to create an empty version
	Init() Foldable
	// there needs to be a way to combine an Item and a Foldable
	Append(item Item) Foldable
	// the generic functions below can return foldables, so we need a conversion method
	AsItem() Item
}

// there's a few things that can be defined with just a (left) fold
// the interface Foldable *cannot* be the receiver of the function, but that just shows that it can work with any type

// Map applies a function to each item inside the foldable
func Map(foldable Foldable, mapFunc func(Item) Item) Foldable {
	result := foldable.Foldl(foldable.Init().AsItem(), func(result, next Item) Item {
		return result.AsFoldable().
			Append(mapFunc(next)).
			AsItem()
	})
	return result.AsFoldable()
}

// Filter returns all the items which pass the filter func
func Filter(foldable Foldable, filterFunc func(Item) bool) Foldable {
	result := foldable.Foldl(foldable.Init().AsItem(), func(result, next Item) Item {
		if filterFunc(next) {
			return result.AsFoldable().
				Append(next).
				AsItem()
		}
		return result
	})
	return result.AsFoldable()
}

// IntItem is being used here, but every Foldable will use it for length

// Length returns the number of items contained in a foldable
func Length(foldable Foldable) int {
	count := IntItem{Value: 0}
	result := foldable.Foldl(count, func(result, next Item) Item {
		return IntItem{Value: result.(IntItem).Value + 1}
	})
	return result.(IntItem).Value
}

// some generic functions operate on boolean values, so we need to define a BoolItem type.
// because Item can be a foldable, we also need to define a BoolFoldable as part of this

// All returns true if all items pass the filterFunc
func All(foldable Foldable, filterFunc func(Item) bool) bool {
	result := foldable.Foldl(BoolItem{Value: true}, func(result, next Item) Item {
		return BoolItem{Value: result.(BoolItem).Value && filterFunc(next)}
	})
	return result.(BoolItem).Value
}

// Any returns true if any of the items pass the filterFunc
func Any(foldable Foldable, filterFunc func(Item) bool) bool {
	result := foldable.Foldl(BoolItem{Value: false}, func(result, next Item) Item {
		return BoolItem{Value: (result.(BoolItem).Value || filterFunc(next))}
	})
	return result.(BoolItem).Value
}

// Concat concatenates the parameters
func Concat(a, b Foldable) Foldable {
	result := b.Foldl(a.AsItem(), func(result, next Item) Item {
		return result.AsFoldable().Append(next).AsItem()
	})
	return result.AsFoldable()
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
