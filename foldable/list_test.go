package foldable

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestListFoldableAppend(t *testing.T) {
	expected := List{1, 2, 3, 4, 5}
	got := List{1, 2, 3, 4}.Append(5).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListMap(t *testing.T) {
	expected := List{0, 2, 4}
	got := Map(
		List{0, 1, 2},
		func(x T) T { return x.(int) * 2 }).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListFilter(t *testing.T) {
	expected := List{-1, -30}
	got := Filter(
		List{0, -1, 1, 2, -30},
		func(x T) bool { return x.(int) < 0 }).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListLength(t *testing.T) {
	length := Length(List{0, -1, 1, 2, -30})
	if length != 5 {
		t.Errorf("result == %v expected %v", length, 5)
	}
}
func TestListAll(t *testing.T) {
	allNegative := All(
		List{0, -1, 1, 2, -30},
		func(x T) bool { return x.(int) < 0 })
	if allNegative {
		t.Errorf("result == %v expected %v", allNegative, false)
	}
}

func TestListAny(t *testing.T) {
	anyNegative := Any(
		List{0, -1, 1, 2, -30},
		func(x T) bool { return x.(int) < 0 })
	if !anyNegative {
		t.Errorf("result == %v expected %v", anyNegative, true)
	}
}

func TestListConcat(t *testing.T) {
	expected := List{1, 2, 3, 4, 5, 6}
	got := Concat(List{1, 2, 3}, List{4, 5, 6}).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListTake(t *testing.T) {
	expected := List{1, 2, 3}
	got := Take(List{1, 2, 3, 4, 5, 6}, 3).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListDrop(t *testing.T) {
	expected := List{4, 5, 6}
	got := Drop(List{1, 2, 3, 4, 5, 6}, 3).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListParMap(t *testing.T) {
	expected := List{2, 4, 6}
	got := ParMap(
		List{1, 2, 3},
		func(x T) T {
			time.Sleep(2 * time.Second)
			fmt.Println("processing Item", x)
			return x.(int) * 2
		}).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListPartition(t *testing.T) {
	expectedPass := List{-1, -30}
	expectedFail := List{0, 1, 2}
	pass, fail := Partition(
		List{0, -1, 1, 2, -30},
		func(x T) bool { return x.(int) < 0 })
	if !reflect.DeepEqual(expectedPass, pass.(List)) {
		t.Errorf("result == %v expected %v", pass, expectedPass)
	}
	if !reflect.DeepEqual(expectedFail, fail.(List)) {
		t.Errorf("result == %v expected %v", fail, expectedFail)
	}
}

func TestListZip(t *testing.T) {
	a := List{1, 2, 3}
	b := List{4, 5, 6, 7}
	expected := List{Pair{1, 4}, Pair{2, 5}, Pair{3, 6}}
	got := Zip(a, b).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
	// length mismatch in other direction
	expected = List{Pair{4, 1}, Pair{5, 2}, Pair{6, 3}}
	got = Zip(b, a).(List)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListUnzip(t *testing.T) {
	a := List{1, 2, 3}
	b := List{4, 5, 6}
	gotA, gotB := Unzip(Zip(a, b))
	if !reflect.DeepEqual(a, gotA) {
		t.Errorf("result == %v expected %v", gotA, a)
	}
	if !reflect.DeepEqual(b, gotB) {
		t.Errorf("result == %v expected %v", gotB, b)
	}
}
