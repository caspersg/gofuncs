package foldable

import (
	"reflect"
	"testing"
)

func TestHashAndListMap(t *testing.T) {
	in := Hash{Value: map[string]T{"a": 1, "b": 2, "c": 3}}
	expected := []T{3, 5, 7}
	got := Map(
		in.Foldl(List{}.Init(), func(result, next T) T {
			return result.(Foldable).Append(next.(HashEntry).Value.(int) * 2)
		}).(Foldable),
		func(x T) T { return x.(int) + 1 }).(List).Value
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestHashAndListMapToType(t *testing.T) {
	in := Hash{Value: map[string]T{"a": 1, "b": 2, "c": 3}}
	expected := []T{3, 5, 7}
	got := Map(
		MapToType(List{}, in, func(x T) T { return x.(HashEntry).Value.(int) * 2 }),
		func(x T) T { return x.(int) + 1 }).(List).Value
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}
