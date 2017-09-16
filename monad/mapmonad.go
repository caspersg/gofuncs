package monad

type K interface{}
type V interface{}

// A key value pair
type Entry struct {
	Key   K
	Value V
}

type MapMonad struct {
	Map map[K]V
}

func (mapMonad MapMonad) Unit(x A) Monad {
	e := x.(Entry)
	return MapMonad{Map: map[K]V{e.Key: e.Value}}
}

func (mapMonad MapMonad) FlatMap(f func(A) Monad) Monad {
	results := map[K]V{}
	for k, v := range mapMonad.Map {
		for k2, y2 := range f(Entry{Key: k, Value: v}).(MapMonad).Map {
			results[k2] = y2
		}
	}

	return MapMonad{Map: results}
}
