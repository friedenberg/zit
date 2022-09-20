package gattung

import "fmt"

type ErrWrongType struct {
	ExpectedType, ActualType Gattung
}

func (e ErrWrongType) Error() string {
	return fmt.Sprintf("expected zk_types %s but got %s", e.ExpectedType, e.ActualType)
}
