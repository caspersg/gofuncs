package monad

import (
	"reflect"
	"testing"
)

func TestMapMonad(t *testing.T) {
	mapFunc := func(x A) Monad {
		e := x.(Entry)
		return MapMonad{}.Unit(Entry{Key: e.Key, Value: (e.Value.(int) * 2)})
	}
	in := map[K]V{"a": 1, "b": 2, "c": 3}
	expected := map[K]V{"a": 2, "b": 4, "c": 6}
	got := MapMonad{Map: in}.FlatMap(mapFunc).(MapMonad)
	if !reflect.DeepEqual(got.Map, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}
