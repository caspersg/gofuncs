package funcs

// bool is required for functions which return a bool
type BoolItem struct {
	Value bool
}

func (boolItem BoolItem) AsFoldable() Foldable {
	return BoolFoldable{List: []bool{boolItem.Value}}
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

func (boolFoldable BoolFoldable) AsItem() Item {
	return boolFoldable
}

func (boolFoldable BoolFoldable) AsFoldable() Foldable {
	return boolFoldable
}
