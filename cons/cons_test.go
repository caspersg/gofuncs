package cons

import (
	"reflect"
	"strings"
	"testing"
)

func TestCar(t *testing.T) {
	got := car(cons(1, nil))
	if !reflect.DeepEqual(got, 1) {
		t.Errorf("result == %v", got)
	}
	got = car(cons(1, cons(2, nil)))
	if !reflect.DeepEqual(got, 1) {
		t.Errorf("result == %v", got)
	}
	got = car(cdr(cons(1, cons(2, nil))))
	if !reflect.DeepEqual(got, 2) {
		t.Errorf("result == %v", got)
	}
}

func TestCdr(t *testing.T) {
	got := cdr(cons(1, nil))
	if !reflect.DeepEqual(got, nil) {
		t.Errorf("result == %v", got)
	}
	got = cdr(cdr(cons(1, cons(2, nil))))
	if !reflect.DeepEqual(got, nil) {
		t.Errorf("result == %v", got)
	}
}

func TestFoldl(t *testing.T) {
	l := cons(1, cons(2, cons(3, nil)))
	add := func(x, y Any) Any {
		return x.(int) + y.(int)
	}
	got := foldl(add, 0, l)
	if !reflect.DeepEqual(got, 6) {
		t.Errorf("result == %v", got)
	}
	mult := func(x, y Any) Any {
		return x.(int) * y.(int)
	}
	got = foldl(mult, 1, l)
	if !reflect.DeepEqual(got, 6) {
		t.Errorf("result == %v", got)
	}
}

func TestFilter(t *testing.T) {
	l := cons(1, cons(2, cons(3, nil)))
	even := func(x Any) bool {
		return x.(int)%2 == 0
	}
	got := filter(even, l)
	if !reflect.DeepEqual(toString(got), "cons(2, nil)") {
		t.Errorf("result == %v", toString(got))
	}
	isThree := func(x Any) bool {
		return x.(int) == 3
	}
	got = filter(isThree, l)
	if !reflect.DeepEqual(toString(got), "cons(3, nil)") {
		t.Errorf("result == %v", toString(got))
	}

	s := cons("afds", cons("dfs", cons("ab", nil)))
	startsWithA := func(x Any) bool {
		return strings.HasPrefix(x.(string), "a")
	}
	got = filter(startsWithA, s)
	if !reflect.DeepEqual(toString(got), "cons(ab, cons(afds, nil))") {
		t.Errorf("result == %v", toString(got))
	}
}
