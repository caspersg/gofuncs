package cons

// haskell foldl
// foldl :: (a -> b -> b) -> b -> [a] -> b
// foldl f z []     = z
// foldl f z (x:xs) = foldl f (f z x) xs
func foldl(f func(Any, Any) Any, z, l Any) Any {
	if l == nil {
		return z
	}
	// recursively apply 'f' to the first and the rest of the items in the list
	return foldl(f, f(z, car(l)), cdr(l))
}

// haskell foldr
// foldr :: (a -> b -> b) -> b -> [a] -> b
// foldr f z []     = z
// foldr f z (x:xs) = f x (foldr f z xs)
func foldr(f func(Any, Any) Any, z, l Any) Any {
	if l == nil {
		return z
	}
	// applies 'f' recursively without tail calls
	return f(car(l), foldr(f, z, cdr(l)))
}

// filter p xs = foldr (\x xs -> if p x then x : xs else xs) [] xs
func filter(p func(Any) bool, l Any) Any {
	return foldr(func(x Any, xs Any) Any {
		if p(x) {
			return cons(x, xs)
		}
		return xs
	}, nil, l)
}

// filter' p xs = foldl (\xs x -> if p x then x : xs else xs) [] xs
// reverses order
func filterl(p func(Any) bool, l Any) Any {
	return foldl(func(xs Any, x Any) Any {
		if p(x) {
			return cons(x, xs)
		}
		return xs
	}, nil, l)
}
