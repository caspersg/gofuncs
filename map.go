package funcs

// MapInt : functional map over int array
func MapInt(f func(int) int, a []int) []int {
	if f == nil {
		return nil
	}
	var res = make([]int, len(a))
	for i, x := range a {
		res[i] = f(x)
	}
	return res
}

// Fold functional fold over int array
func Fold(init int, f func(int, int) int, a []int) int {
	var res = init
	for _, x := range a {
		res = f(res, x)
	}
	return res
}

// SliceEqual compares slices
func SliceEqual(a []int, b []int) bool {
	if len(a) == len(b) {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	return false
}
