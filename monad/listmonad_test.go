package monad

import (
	"reflect"
	"strings"
	"testing"
)

func TestIntListMonad(t *testing.T) {
	mapFunc := func(x A) Monad { return ListMonad{}.Unit(x.(int) * 2) }
	in := []A{0, 1, 2}
	expected := []A{0, 2, 4}
	got := ListMonad{List: in}.FlatMap(mapFunc)
	if !reflect.DeepEqual(got.(ListMonad).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestStringListMonad(t *testing.T) {
	mapFunc := func(x A) Monad { return ListMonad{}.Unit(strings.ToUpper(x.(string))) }
	in := []A{"a", "abc", "aBcD"}
	expected := []A{"A", "ABC", "ABCD"}
	got := ListMonad{List: in}.FlatMap(mapFunc)
	if !reflect.DeepEqual(got.(ListMonad).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntAndStringListMonad(t *testing.T) {
	mapFunc := func(x A) Monad { return ListMonad{}.Unit(len(x.(string))) }
	in := []A{"a", "abc", "aBcD"}
	expected := []A{1, 3, 4}
	got := ListMonad{List: in}.FlatMap(mapFunc)
	if !reflect.DeepEqual(got.(ListMonad).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}
