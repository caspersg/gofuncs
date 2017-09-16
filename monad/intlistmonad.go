package monad

// T int

// IntListMonad M[T]
type IntListMonad struct {
	List []int
}

// Unit A -> M[A]
func (intListMonad IntListMonad) Unit(x T) Monad {
	return IntListMonad{List: []int{x.(int)}}
}

// FlatMap M[A] -> (A -> M[B]) -> M[B]
func (intListMonad IntListMonad) FlatMap(f func(T) Monad) Monad {
	results := []int{}
	for _, x := range intListMonad.List {
		for _, y := range f(x).(IntListMonad).List {
			results = append(results, y)
		}
	}
	return IntListMonad{List: results}
}

// Map M[A] -> (A->B) -> M[B]
// func (intListMonad IntListMonad) Map(f func(T) Monad) Monad {
// 	results := []Monad
// 	for _, x := range intListMonad {
// 		f(x)
// 	}
// }

// // flatten M[M[B]] -> M[B]
// func Flatten(listOfLists []Monad) Monad {
// 	results := []int{}
// 	for _, x := range listOfLists {
// 		for _, y := range x.(IntListMonad).List {
// 			results = append(results, y)
// 		}
// 	}
// 	return IntListMonad{List: results}
// }
