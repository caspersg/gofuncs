package foldable

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestHashAppend(t *testing.T) {
	got := Hash{"b": 2}.Append(HashEntry{Key: "a", Value: 5}).(Hash)
	expected := Hash{"b": 2, "a": 5}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestHashMap(t *testing.T) {
	expected := Hash{"a": 2, "b": 4, "c": 6}
	got := Map(Hash{"a": 1, "b": 2, "c": 3}, func(x T) T {
		return HashEntry{Key: x.(HashEntry).Key, Value: x.(HashEntry).Value.(int) * 2}
	}).(Hash)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestHashFilter(t *testing.T) {
	expected := Hash{"b": -2, "d": -30}
	got := Filter(
		Hash{"a": 1, "b": -2, "c": 3, "d": -30},
		func(x T) bool { return x.(HashEntry).Value.(int) < 0 }).(Hash)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestHashLength(t *testing.T) {
	length := Length(Hash{"a": 1, "b": -2, "c": 3, "d": -30})
	if length != 4 {
		t.Errorf("result == %v expected %v", length, 5)
	}
}
func TestHashAll(t *testing.T) {
	allNegative := All(
		Hash{"a": 1, "b": -2, "c": 3, "d": -30},
		func(x T) bool { return x.(HashEntry).Value.(int) < 0 })
	if allNegative {
		t.Errorf("result == %v expected %v", allNegative, false)
	}
}

func TestHashAny(t *testing.T) {
	anyB := Any(
		Hash{"a": 1, "b": -2, "c": 3, "d": -30},
		func(x T) bool { return x.(HashEntry).Key == "b" })
	if !anyB {
		t.Errorf("result == %v expected %v", anyB, true)
	}
}

func TestHashConcat(t *testing.T) {
	expected := Hash{"a": 1, "b": -2, "c": 3, "d": -30}
	result := Concat(Hash{"a": 1, "b": -2}, Hash{"c": 3, "d": -30})
	if !reflect.DeepEqual(result.(Hash), expected) {
		t.Errorf("result == %v expected %v", result, expected)
	}
}

func TestHashTake(t *testing.T) {
	expected := Hash{"a": 1, "b": -2, "c": 3}
	result := Take(Hash{"a": 1, "b": -2, "c": 3, "d": 4, "e": 5, "f": 6}, 3)
	if !reflect.DeepEqual(result.(Hash), expected) {
		t.Errorf("result == %v expected %v", result, expected)
	}
}

func TestHashDrop(t *testing.T) {
	expected := Hash{"d": 4, "e": 5, "f": 6}
	result := Drop(Hash{"a": 1, "b": -2, "c": 3, "d": 4, "e": 5, "f": 6}, 3)
	if !reflect.DeepEqual(result.(Hash), expected) {
		t.Errorf("result == %v expected %v", result, expected)
	}
}

func TestHashParMap(t *testing.T) {
	expected := Hash{"a": 2, "b": 4, "c": 6}
	got := ParMap(
		Hash{"a": 1, "b": 2, "c": 3},
		func(x T) T {
			time.Sleep(2 * time.Second)
			fmt.Println("processing Item", x)
			return HashEntry{Key: x.(HashEntry).Key, Value: x.(HashEntry).Value.(int) * 2}
		})
	if !reflect.DeepEqual(got.(Hash), expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}
