package listable

import (
	"fmt"
	// "reflect"
	"testing"
	"time"
)

func listContentsEqual(expected []int, got []T) bool {
	if len(expected) != len(got) {
		return false
	}
	for i, x := range expected {
		if got[i] != x {
			return false
		}
	}
	return true
}

func TestIntFoldableMap(t *testing.T) {
	mapFunc := func(x T) T { return x.(int) * 2 }
	// no auto-boxing of base types into struct, so the below line cannot work
	// got := Map([]int{0, 1, 2}.([]T), mapFunc)
	// therefore the type of the input MUST be []T
	in := []T{0, 1, 2}
	expected := []int{0, 2, 4}
	got := Map(in, mapFunc)
	if !listContentsEqual(expected, got) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableFilter(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []T{0, -1, 1, 2, -30}
	expected := []int{-1, -30}
	got := Filter(in, isNegative)
	if !listContentsEqual(expected, got) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableLength(t *testing.T) {
	in := []T{0, -1, 1, 2, -30}
	length := Length(in)
	if length != 5 {
		t.Errorf("result == %d expected %d", length, 5)
	}
}
func TestIntFoldableAll(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []T{0, -1, 1, 2, -30}
	allNegative := All(in, isNegative)
	if allNegative {
		t.Errorf("result == %d expected %d", allNegative, false)
	}
}

func TestIntFoldableAny(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []T{0, -1, 1, 2, -30}
	anyNegative := Any(in, isNegative)
	if !anyNegative {
		t.Errorf("result == %d expected %d", anyNegative, true)
	}
}

func TestIntFoldableConcat(t *testing.T) {
	a := []T{1, 2, 3}
	b := []T{4, 5, 6}
	expected := []int{1, 2, 3, 4, 5, 6}
	got := Concat(a, b)
	if !listContentsEqual(expected, got) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

// func TestIntFoldableTake(t *testing.T) {
// 	x := []int{1, 2, 3, 4, 5, 6}
// 	expected := []int{1, 2, 3}
// 	result := Take(IntFoldable{List: x}, 3)
// 	if !reflect.DeepEqual(result.(IntFoldable).List, expected) {
// 		t.Errorf("result == %d expected %d", result, expected)
// 	}
// }

// func TestIntFoldableDrop(t *testing.T) {
// 	x := []int{1, 2, 3, 4, 5, 6}
// 	expected := []int{4, 5, 6}
// 	result := Drop(IntFoldable{List: x}, 3)
// 	if !reflect.DeepEqual(result.(IntFoldable).List, expected) {
// 		t.Errorf("result == %d expected %d", result, expected)
// 	}
// }

func TestIntFoldableParMap(t *testing.T) {
	mapFunc := func(x T) T {
		time.Sleep(2 * time.Second)
		fmt.Println("processing Item", x)
		return x.(int) * 2
	}
	in := []T{1, 2, 3}
	expected := []int{2, 4, 6}
	got := ParMap(in, mapFunc)
	if !listContentsEqual(expected, got) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}
