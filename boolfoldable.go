package funcs

type BoolItem struct {
	Value bool
}

type BoolFoldable struct {
	List []bool
}

func (boolFoldable BoolFoldable) Foldl(init Item, foldFunc func(result, next Item) Item) Item {
	result := init
	for _, x := range boolFoldable.List {
		result = foldFunc(result, BoolItem{Value: x})
	}
	return result
}

func (boolFoldable BoolFoldable) Init() Foldable {
	return BoolFoldable{}
}

func (boolFoldable BoolFoldable) Append(item Item) Foldable {
	return BoolFoldable{List: append(boolFoldable.List, item.(BoolItem).Value)}
}
