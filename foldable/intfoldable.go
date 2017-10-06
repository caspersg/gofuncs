package foldable

type IntFoldable []int

func (intFoldable IntFoldable) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	for _, x := range intFoldable {
		result = foldFunc(result, x)
	}
	return result
}

func (intFoldable IntFoldable) Init() Foldable {
	return IntFoldable{}
}

func (intFoldable IntFoldable) Append(item T) Foldable {
	return append(intFoldable, item.(int))
}
