package cons

import "fmt"

//(define (cons x y)
//  (lambda (m) (m x y)))
//(define (car z)
//  (z (lambda (p q) p)))
//(define (cdr z)
//  (z (lambda (p q) q)))

type Any interface{}

func cons(x, y Any) Any {
	return func(m func(Any, Any) Any) Any {
		return m(x, y)
	}
}

func car(z Any) Any {
	zf := z.(func(func(Any, Any) Any) Any)
	return zf(func(p, q Any) Any {
		return p
	})
}

func cdr(z Any) Any {
	zf := z.(func(func(Any, Any) Any) Any)
	return zf(func(p, q Any) Any {
		return q
	})
}

func toString(l Any) string {
	if l == nil {
		return "nil"
	}
	return fmt.Sprintf("cons(%v, %v)", car(l), toString(cdr(l)))
}
