package foldable

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestListFoldableAppend(t *testing.T) {
	anItem := 5
	aList := List{Value: []T{1, 2, 3, 4}}
	expected := []T{1, 2, 3, 4, 5}
	got := aList.Append(anItem).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListMap(t *testing.T) {
	mapFunc := func(x T) T { return x.(int) * 2 }
	in := []T{0, 1, 2}
	expected := []T{0, 2, 4}
	got := Map(List{Value: in}, mapFunc).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListFilter(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []T{0, -1, 1, 2, -30}
	expected := []T{-1, -30}
	got := Filter(List{Value: in}, isNegative).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListLength(t *testing.T) {
	in := []T{0, -1, 1, 2, -30}
	length := Length(List{Value: in})
	if length != 5 {
		t.Errorf("result == %v expected %v", length, 5)
	}
}
func TestListAll(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []T{0, -1, 1, 2, -30}
	allNegative := All(List{Value: in}, isNegative)
	if allNegative {
		t.Errorf("result == %v expected %v", allNegative, false)
	}
}

func TestListAny(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []T{0, -1, 1, 2, -30}
	anyNegative := Any(List{Value: in}, isNegative)
	if !anyNegative {
		t.Errorf("result == %v expected %v", anyNegative, true)
	}
}

func TestListConcat(t *testing.T) {
	a := []T{1, 2, 3}
	b := []T{4, 5, 6}
	expected := []T{1, 2, 3, 4, 5, 6}
	got := Concat(List{Value: a}, List{Value: b}).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListTake(t *testing.T) {
	x := []T{1, 2, 3, 4, 5, 6}
	expected := []T{1, 2, 3}
	got := Take(List{Value: x}, 3).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListDrop(t *testing.T) {
	x := []T{1, 2, 3, 4, 5, 6}
	expected := []T{4, 5, 6}
	got := Drop(List{Value: x}, 3).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListParMap(t *testing.T) {
	mapFunc := func(x T) T {
		time.Sleep(2 * time.Second)
		fmt.Println("processing Item", x)
		return x.(int) * 2
	}
	in := []T{1, 2, 3}
	expected := []T{2, 4, 6}
	got := ParMap(List{Value: in}, mapFunc).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListPartition(t *testing.T) {
	isNegative := func(x T) bool { return x.(int) < 0 }
	in := []T{0, -1, 1, 2, -30}
	expectedPass := []T{-1, -30}
	expectedFail := []T{0, 1, 2}
	pass, fail := Partition(List{Value: in}, isNegative)
	if !reflect.DeepEqual(expectedPass, pass.(List).Value) {
		t.Errorf("result == %v expected %v", pass, expectedPass)
	}
	if !reflect.DeepEqual(expectedFail, fail.(List).Value) {
		t.Errorf("result == %v expected %v", fail, expectedFail)
	}
}

func TestListZip(t *testing.T) {
	a := []T{1, 2, 3}
	b := []T{4, 5, 6, 7}
	expected := []T{Pair{1, 4}, Pair{2, 5}, Pair{3, 6}}
	got := Zip(List{Value: a}, List{Value: b}).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}

	expected = []T{Pair{4, 1}, Pair{5, 2}, Pair{6, 3}}
	got = Zip(List{Value: b}, List{Value: a}).(List).Value
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}
