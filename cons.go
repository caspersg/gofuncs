package funcs

type F1 func(F2) Any
type F2 func(Any, Any) Any
type Any interface{}

// Cons construct list using only functions
func Cons(a Any, b Any) F1 {
	return func(f F2) Any { return f(a, b) }
}

func Car(z F1) Any {
	return z(func(a Any, b Any) Any { return a })
}

func Cdr(z F1) Any {
	return z(func(a Any, b Any) Any { return b })
}
