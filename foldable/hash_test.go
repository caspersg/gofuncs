package foldable

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestHashAppend(t *testing.T) {
	anItem := HashEntry{Key: "a", Value: 5}
	in := Hash{Value: map[string]T{"b": 2}}
	got := in.Append(anItem).(Hash).Value
	expected := map[string]T{"b": 2, "a": 5}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestHashMap(t *testing.T) {
	mapFunc := func(x T) T {
		return HashEntry{Key: x.(HashEntry).Key, Value: x.(HashEntry).Value.(int) * 2}
	}
	in := Hash{Value: map[string]T{"a": 1, "b": 2, "c": 3}}
	expected := map[string]T{"a": 2, "b": 4, "c": 6}
	got := Map(in, mapFunc).(Hash).Value
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestHashFilter(t *testing.T) {
	isNegative := func(x T) bool { return x.(HashEntry).Value.(int) < 0 }
	in := Hash{Value: map[string]T{"a": 1, "b": -2, "c": 3, "d": -30}}
	expected := map[string]T{"b": -2, "d": -30}
	got := Filter(in, isNegative).(Hash).Value
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestHashLength(t *testing.T) {
	in := Hash{Value: map[string]T{"a": 1, "b": -2, "c": 3, "d": -30}}
	length := Length(in)
	if length != 4 {
		t.Errorf("result == %v expected %v", length, 5)
	}
}
func TestHashAll(t *testing.T) {
	isNegative := func(x T) bool { return x.(HashEntry).Value.(int) < 0 }
	in := Hash{Value: map[string]T{"a": 1, "b": -2, "c": 3, "d": -30}}
	allNegative := All(in, isNegative)
	if allNegative {
		t.Errorf("result == %v expected %v", allNegative, false)
	}
}

func TestHashAny(t *testing.T) {
	isB := func(x T) bool { return x.(HashEntry).Key == "b" }
	in := Hash{Value: map[string]T{"a": 1, "b": -2, "c": 3, "d": -30}}
	anyB := Any(in, isB)
	if !anyB {
		t.Errorf("result == %v expected %v", anyB, true)
	}
}

func TestHashConcat(t *testing.T) {
	a := Hash{Value: map[string]T{"a": 1, "b": -2}}
	b := Hash{Value: map[string]T{"c": 3, "d": -30}}
	expected := map[string]T{"a": 1, "b": -2, "c": 3, "d": -30}
	result := Concat(a, b)
	if !reflect.DeepEqual(result.(Hash).Value, expected) {
		t.Errorf("result == %v expected %v", result, expected)
	}
}

func TestHashTake(t *testing.T) {
	a := Hash{Value: map[string]T{"a": 1, "b": -2, "c": 3, "d": 4, "e": 5, "f": 6}}
	expected := map[string]T{"a": 1, "b": -2, "c": 3}
	result := Take(a, 3)
	if !reflect.DeepEqual(result.(Hash).Value, expected) {
		t.Errorf("result == %v expected %v", result, expected)
	}
}

func TestHashDrop(t *testing.T) {
	a := Hash{Value: map[string]T{"a": 1, "b": -2, "c": 3, "d": 4, "e": 5, "f": 6}}
	expected := map[string]T{"d": 4, "e": 5, "f": 6}
	result := Drop(a, 3)
	if !reflect.DeepEqual(result.(Hash).Value, expected) {
		t.Errorf("result == %v expected %v", result, expected)
	}
}

func TestHashParMap(t *testing.T) {
	mapFunc := func(x T) T {
		time.Sleep(2 * time.Second)
		fmt.Println("processing Item", x)
		return HashEntry{Key: x.(HashEntry).Key, Value: x.(HashEntry).Value.(int) * 2}
	}
	in := Hash{Value: map[string]T{"a": 1, "b": 2, "c": 3}}
	expected := map[string]T{"a": 2, "b": 4, "c": 6}
	got := ParMap(in, mapFunc)
	if !reflect.DeepEqual(got.(Hash).Value, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}
