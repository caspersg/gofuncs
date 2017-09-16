package monad

type ListMonad struct {
	List []A
}

func (listMonad ListMonad) Unit(x A) Monad {
	return ListMonad{List: []A{x}}
}

func (listMonad ListMonad) FlatMap(f func(A) Monad) Monad {
	results := []A{}
	for _, x := range listMonad.List {
		for _, y := range f(x).(ListMonad).List {
			results = append(results, y)
		}
	}

	return ListMonad{List: results}
}
