package cons

type Any interface{}

//(define (cons x y)
//  (lambda (m) (m x y)))
func cons(x, y Any) Any {
	// x and y are passed in on construction
	return func(m func(Any, Any) Any) Any {
		// x an y are captured in the returned closure
		// the closure is now how these values are stored

		// the returned function will apply any function 'm' to the captured values
		return m(x, y)
	}
}

//(define (car z)
//  (z (lambda (p q) p)))
func car(z Any) Any {
	zf := z.(func(func(Any, Any) Any) Any)
	// when we apply zf, it will in turn apply the funcion we pass it to its captured values
	return zf(func(p, q Any) Any {
		// return the first captured valued for this 'cons' (part of the list)
		return p
	})
}

//(define (cdr z)
//  (z (lambda (p q) q)))
func cdr(z Any) Any {
	zf := z.(func(func(Any, Any) Any) Any)
	return zf(func(p, q Any) Any {
		// same as 'car' but it returns the second captured valued
		return q
	})
}
