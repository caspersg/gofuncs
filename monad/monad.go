package monad

type T interface{}

// Monad this is an alternative to Foldable
type Monad interface {
	FlatMap(f func(T) Monad) Monad
	Unit(x T) Monad
}
