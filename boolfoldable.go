package funcs

type BoolFoldable struct {
	List []bool
}

func (boolFoldable BoolFoldable) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	for _, x := range boolFoldable.List {
		result = foldFunc(result, x)
	}
	return result
}

func (boolFoldable BoolFoldable) Init() Foldable {
	return BoolFoldable{}
}

func (boolFoldable BoolFoldable) Append(item T) Foldable {
	return BoolFoldable{List: append(boolFoldable.List, item.(bool))}
}
