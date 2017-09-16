package monad

import (
	"reflect"
	"testing"
)

func TestIntListMonad(t *testing.T) {
	mapFunc := func(x T) Monad { return IntListMonad{}.Unit(x.(int) * 2) }
	in := []int{0, 1, 2}
	expected := []int{0, 2, 4}
	got := IntListMonad{List: in}.FlatMap(mapFunc)
	if !reflect.DeepEqual(got.(IntListMonad).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}
