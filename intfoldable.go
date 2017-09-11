package funcs

// lets use int as an example implementation

type IntFoldable struct {
	List []int
}

func (intFoldable IntFoldable) Foldl(init Item, foldFunc func(result, next Item) Item) Item {
	result := init
	for _, x := range intFoldable.List {
		result = foldFunc(result, x)
	}
	return result
}

func (intFoldable IntFoldable) Init() Foldable {
	return IntFoldable{}
}

func (intFoldable IntFoldable) Append(item Item) Foldable {
	return IntFoldable{List: append(intFoldable.List, item.(int))}
}
