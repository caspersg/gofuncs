package monad

import "fmt"

// T string

// StringListMonad M[T]
type StringListMonad struct {
	List []string
}

// Unit A -> M[A]
func (stringListMonad StringListMonad) Unit(x T) Monad {
	return StringListMonad{List: []string{x.(string)}}
}

// FlatMap M[A] -> (A -> M[B]) -> M[B]
func (stringListMonad StringListMonad) FlatMap(f func(T) Monad) Monad {
	results := []string{}
	for _, x := range stringListMonad.List {
		fmt.Println("x", x)
		for _, y := range f(x).(StringListMonad).List {
			fmt.Println("y", y)
			results = append(results, y)
		}
	}
	return StringListMonad{List: results}
}
