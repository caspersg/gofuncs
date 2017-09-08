package funcs

// lets use int as an example implementation

type IntItem struct {
	Value int
}

func (intItem IntItem) AsFoldable() Foldable {
	return IntFoldable{List: []int{intItem.Value}}
}

type IntFoldable struct {
	List []int
}

func (intList IntFoldable) Foldl(init Item, foldFunc func(result, next Item) Item) Item {
	result := init
	for _, x := range intList.List {
		result = foldFunc(result, IntItem{Value: x})
	}
	return result
}

func (intList IntFoldable) Init() Foldable {
	return IntFoldable{}
}

func (intList IntFoldable) Append(item Item) Foldable {
	return IntFoldable{List: append(intList.List, item.(IntItem).Value)}
}

func (intList IntFoldable) AsItem() Item {
	return intList
}

func (intList IntFoldable) AsFoldable() Foldable {
	return intList
}
