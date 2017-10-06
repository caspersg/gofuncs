package foldable

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestIntFoldableAppend(t *testing.T) {
	appendResult := IntFoldable{1, 2, 3, 4}.Append(5).(IntFoldable)
	expected := IntFoldable{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(appendResult, expected) {
		t.Errorf("result == %d expected %d", appendResult, expected)
	}
}

func TestIntFoldableMap(t *testing.T) {
	expected := IntFoldable{0, 2, 4}
	got := Map(IntFoldable{0, 1, 2}, func(x T) T { return x.(int) * 2 })
	if !reflect.DeepEqual(got.(IntFoldable), expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableFilter(t *testing.T) {
	expected := IntFoldable{-1, -30}
	got := Filter(IntFoldable{0, -1, 1, 2, -30}, func(x T) bool { return x.(int) < 0 })
	if !reflect.DeepEqual(got.(IntFoldable), expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableLength(t *testing.T) {
	length := Length(IntFoldable{0, -1, 1, 2, -30})
	if length != 5 {
		t.Errorf("result == %d expected %d", length, 5)
	}
}
func TestIntFoldableAll(t *testing.T) {
	allNegative := All(IntFoldable{0, -1, 1, 2, -30}, func(x T) bool { return x.(int) < 0 })
	if allNegative {
		t.Errorf("result == %d expected %d", allNegative, false)
	}
}

func TestIntFoldableAny(t *testing.T) {
	anyNegative := Any(IntFoldable{0, -1, 1, 2, -30}, func(x T) bool { return x.(int) < 0 })
	if !anyNegative {
		t.Errorf("result == %d expected %d", anyNegative, true)
	}
}

func TestIntFoldableConcat(t *testing.T) {
	expected := IntFoldable{1, 2, 3, 4, 5, 6}
	result := Concat(IntFoldable{1, 2, 3}, IntFoldable{4, 5, 6})
	if !reflect.DeepEqual(result.(IntFoldable), expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableTake(t *testing.T) {
	expected := IntFoldable{1, 2, 3}
	result := Take(IntFoldable{1, 2, 3, 4, 5, 6}, 3)
	if !reflect.DeepEqual(result.(IntFoldable), expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableDrop(t *testing.T) {
	expected := IntFoldable{4, 5, 6}
	result := Drop(IntFoldable{1, 2, 3, 4, 5, 6}, 3)
	if !reflect.DeepEqual(result.(IntFoldable), expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableParMap(t *testing.T) {
	expected := IntFoldable{2, 4, 6}
	got := ParMap(IntFoldable{1, 2, 3}, func(x T) T {
		time.Sleep(2 * time.Second)
		fmt.Println("processing Item", x)
		return x.(int) * 2
	})
	if !reflect.DeepEqual(got.(IntFoldable), expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}
