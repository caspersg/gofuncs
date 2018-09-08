package cons

// haskell fold
//foldl f z []     = z
//foldl f z (x:xs) = foldl f (f z x) xs

func foldl(f func(Any, Any) Any, z, l Any) Any {
	if l == nil {
		return z
	}
	return foldl(f, f(z, car(l)), cdr(l))
}

// filter' p xs = foldl (\xs x -> if p x then x : xs else xs) [] xs
// reverses order
func filter(p func(Any) bool, l Any) Any {
	return foldl(func(xs Any, x Any) Any {
		if p(x) {
			return cons(x, xs)
		}
		return xs
	}, nil, l)
}
