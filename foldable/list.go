package foldable

// List is a list of anything
type List struct {
	Value []T
}

func (list List) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	for _, x := range list.Value {
		result = foldFunc(result, x)
	}
	return result
}

func (list List) Init() Foldable {
	return List{}
}

func (list List) Append(item T) Foldable {
	return List{Value: append(list.Value, item)}
}
