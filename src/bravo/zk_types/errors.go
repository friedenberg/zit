package zk_types

import "fmt"

type ErrWrongType struct {
	ExpectedType, ActualType Type
}

func (e ErrWrongType) Error() string {
	return fmt.Sprintf("expected zk_types %s but got %s", e.ExpectedType, e.ActualType)
}
