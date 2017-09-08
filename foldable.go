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
type Item interface {
	AsFoldable() Foldable
}

type Foldable interface {
	Foldl(init Item, f func(result, next Item) Item) Item
	Init() Foldable
	Append(item Item) Foldable
	AsItem() Item
}

// there's a few things that can be defined with just a (left) fold
// the interface Foldable *cannot* be the receiver of the function, but that just shows that it can work with any type
// IntItem is being used here, but every Foldable will use it for length

// Length returns the number of items contained in a foldable
func Length(foldable Foldable) int {
	count := IntItem{Value: 0}
	result := foldable.Foldl(count, func(result, next Item) Item {
		return IntItem{Value: result.(IntItem).Value + 1}
	})
	return result.(IntItem).Value
}

// Map applies a function to each item inside the foldable
func Map(foldable Foldable, mapFunc func(Item) Item) Foldable {
	result := foldable.Foldl(foldable.Init().AsItem(), func(result, next Item) Item {
		f := result.AsFoldable()
		mfr := mapFunc(next)
		ar := f.Append(mfr)
		return ar.AsItem()
		// return result.AsFoldable().Append(mapFunc(next)).AsItem()(Item)
	})
	return result.AsFoldable()
}

// Filter returns all the items which pass the filter func
func Filter(foldable Foldable, filterFunc func(Item) bool) Foldable {
	result := foldable.Foldl(foldable.Init().AsItem(), func(result, next Item) Item {
		if filterFunc(next) {
			return result.AsFoldable().Append(next).AsItem()
		}
		return result
	})
	return result.AsFoldable()
}
