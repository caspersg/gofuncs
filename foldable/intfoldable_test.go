package foldable

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestIntFoldableAppend(t *testing.T) {
	anItem := 5

	aList := IntFoldable{List: []int{1, 2, 3, 4}}
	appendResult := aList.Append(anItem).(IntFoldable).List
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(appendResult, expected) {
		t.Errorf("result == %d expected %d", appendResult, expected)
	}
}

func TestIntFoldableMap(t *testing.T) {
	mapFunc := func(x T) T { return x.(int) * 2 }
	in := []int{0, 1, 2}
	expected := []int{0, 2, 4}
	got := Map(IntFoldable{List: in}, mapFunc)
	if !reflect.DeepEqual(got.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableFilter(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []int{0, -1, 1, 2, -30}
	expected := []int{-1, -30}
	got := Filter(IntFoldable{List: in}, isNegative)
	if !reflect.DeepEqual(got.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableLength(t *testing.T) {
	in := []int{0, -1, 1, 2, -30}
	length := Length(IntFoldable{List: in})
	if length != 5 {
		t.Errorf("result == %d expected %d", length, 5)
	}
}
func TestIntFoldableAll(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []int{0, -1, 1, 2, -30}
	allNegative := All(IntFoldable{List: in}, isNegative)
	if allNegative {
		t.Errorf("result == %d expected %d", allNegative, false)
	}
}

func TestIntFoldableAny(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []int{0, -1, 1, 2, -30}
	anyNegative := Any(IntFoldable{List: in}, isNegative)
	if !anyNegative {
		t.Errorf("result == %d expected %d", anyNegative, true)
	}
}

func TestIntFoldableConcat(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	expected := []int{1, 2, 3, 4, 5, 6}
	result := Concat(IntFoldable{List: a}, IntFoldable{List: b})
	if !reflect.DeepEqual(result.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableTake(t *testing.T) {
	x := []int{1, 2, 3, 4, 5, 6}
	expected := []int{1, 2, 3}
	result := Take(IntFoldable{List: x}, 3)
	if !reflect.DeepEqual(result.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableDrop(t *testing.T) {
	x := []int{1, 2, 3, 4, 5, 6}
	expected := []int{4, 5, 6}
	result := Drop(IntFoldable{List: x}, 3)
	if !reflect.DeepEqual(result.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableParMap(t *testing.T) {
	mapFunc := func(x T) T {
		time.Sleep(2 * time.Second)
		fmt.Println("processing Item", x)
		return x.(int) * 2
	}
	in := []int{1, 2, 3}
	expected := []int{2, 4, 6}
	got := ParMap(IntFoldable{List: in}, mapFunc)
	if !reflect.DeepEqual(got.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}