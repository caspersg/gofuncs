package funcs

import "testing"

func TestCons(t *testing.T) {
	cases := []struct {
		in  Any
		out Any
	}{
		{Car(Cons(10, nil)), 10},
		{Car(Car(Cons(20, Cons(10, nil)))), 20},
		// {Car(Cdr(Cons(1, Cons(2, nil)))), 1},
	}
	for _, c := range cases {
		if c.in != c.out {
			t.Errorf("%d != %d", c.in, c.out)
		}
	}
}
