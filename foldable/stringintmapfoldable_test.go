package foldable

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestStringIntMapFoldableAppend(t *testing.T) {
	anItem := StringIntEntryItem{Key: "a", Value: 5}
	aMap := StringIntMapFoldable{Map: map[string]int{"b": 2}}
	appendResult := aMap.Append(anItem).(StringIntMapFoldable).Map
	expected := map[string]int{"b": 2, "a": 5}
	if !reflect.DeepEqual(appendResult, expected) {
		t.Errorf("result == %v expected %v", appendResult, expected)
	}
}

func TestStringIntMapFoldableMap(t *testing.T) {
	mapFunc := func(x T) T {
		return StringIntEntryItem{Key: x.(StringIntEntryItem).Key, Value: x.(StringIntEntryItem).Value * 2}
	}
	in := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": 2, "c": 3}}
	expected := map[string]int{"a": 2, "b": 4, "c": 6}
	got := Map(in, mapFunc)
	if !reflect.DeepEqual(got.(StringIntMapFoldable).Map, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestStringIntMapFoldableFilter(t *testing.T) {
	isNegative := func(x T) bool { return x.(StringIntEntryItem).Value < 0 }
	in := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": -2, "c": 3, "d": -30}}
	expected := map[string]int{"b": -2, "d": -30}
	got := Filter(in, isNegative)
	if !reflect.DeepEqual(got.(StringIntMapFoldable).Map, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestStringIntMapFoldableLength(t *testing.T) {
	in := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": -2, "c": 3, "d": -30}}
	length := Length(in)
	if length != 4 {
		t.Errorf("result == %d expected %d", length, 5)
	}
}
func TestStringIntMapFoldableAll(t *testing.T) {
	isNegative := func(x T) bool { return x.(StringIntEntryItem).Value < 0 }
	in := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": -2, "c": 3, "d": -30}}
	allNegative := All(in, isNegative)
	if allNegative {
		t.Errorf("result == %d expected %d", allNegative, false)
	}
}

func TestStringIntMapFoldableAny(t *testing.T) {
	isB := func(x T) bool { return x.(StringIntEntryItem).Key == "b" }
	in := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": -2, "c": 3, "d": -30}}
	anyB := Any(in, isB)
	if !anyB {
		t.Errorf("result == %d expected %d", anyB, true)
	}
}

func TestStringIntMapFoldableConcat(t *testing.T) {
	a := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": -2}}
	b := StringIntMapFoldable{Map: map[string]int{"c": 3, "d": -30}}
	expected := map[string]int{"a": 1, "b": -2, "c": 3, "d": -30}
	result := Concat(a, b)
	if !reflect.DeepEqual(result.(StringIntMapFoldable).Map, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestStringIntMapFoldableTake(t *testing.T) {
	a := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": -2, "c": 3, "d": 4, "e": 5, "f": 6}}
	expected := map[string]int{"a": 1, "b": -2, "c": 3}
	result := Take(a, 3)
	if !reflect.DeepEqual(result.(StringIntMapFoldable).Map, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestStringIntMapFoldableDrop(t *testing.T) {
	a := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": -2, "c": 3, "d": 4, "e": 5, "f": 6}}
	expected := map[string]int{"d": 4, "e": 5, "f": 6}
	result := Drop(a, 3)
	if !reflect.DeepEqual(result.(StringIntMapFoldable).Map, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestStringIntMapFoldableParMap(t *testing.T) {
	mapFunc := func(x T) T {
		time.Sleep(2 * time.Second)
		fmt.Println("processing Item", x)
		return StringIntEntryItem{Key: x.(StringIntEntryItem).Key, Value: x.(StringIntEntryItem).Value * 2}
	}
	in := StringIntMapFoldable{Map: map[string]int{"a": 1, "b": 2, "c": 3}}
	expected := map[string]int{"a": 2, "b": 4, "c": 6}
	got := ParMap(in, mapFunc)
	if !reflect.DeepEqual(got.(StringIntMapFoldable).Map, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}
