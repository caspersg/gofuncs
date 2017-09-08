package funcs

import "testing"

func TestMapInt(t *testing.T) {
	cases := []struct {
		inF  func(int) int
		inA  []int
		want []int
	}{
		{nil, nil, nil},
		{nil, []int{0, 1, 2}, []int{}},
		{func(x int) int { return x + 1 }, []int{0, 1, 2}, []int{1, 2, 3}},
		{func(x int) int { return x * 2 }, []int{0, 1, 2}, []int{0, 2, 4}},
	}
	for _, c := range cases {
		got := MapInt(c.inF, c.inA)
		if !SliceEqual(got, c.want) {
			t.Errorf("Map(f, %d) == %d, want %d", c.inA, got, c.want)
		}
	}
}
