package errors

import "fmt"

type wrapped struct {
	outer, inner error
}

func (e wrapped) Error() string {
	return fmt.Sprintf("%s: %s", e.outer, e.inner)
}

func (e wrapped) Unwrap() error {
	return e.inner
}
