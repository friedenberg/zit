package object_id_provider

import "fmt"

type ErrDoesNotExist struct {
	Value string
}

func (e ErrDoesNotExist) Error() string {
	return fmt.Sprintf("object id does not exist: %s", e.Value)
}

func (e ErrDoesNotExist) Is(target error) bool {
	_, ok := target.(ErrDoesNotExist)
	return ok
}

type ErrZettelIdsExhausted struct{}

func (e ErrZettelIdsExhausted) Error() string {
	return "hinweisen exhausted"
}

func (e ErrZettelIdsExhausted) Is(target error) bool {
	_, ok := target.(ErrZettelIdsExhausted)
	return ok
}
