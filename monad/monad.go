package monad

type A interface{}
type B interface{}

// Monad this is an alternative to Foldable
// M[A]
type Monad interface {
	// FlatMap M[A] -> (A -> M[B]) -> M[B]
	FlatMap(f func(A) Monad) Monad
	// Unit A -> M[A]
	Unit(x A) Monad
}
