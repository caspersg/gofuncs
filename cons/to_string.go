package cons

import "fmt"

func prettyPrint(l Any) string {
	if l == nil {
		return "nil"
	}
	return fmt.Sprintf("cons(%v, %v)", car(l), prettyPrint(cdr(l)))
}
