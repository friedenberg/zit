package hinweisen

import "fmt"

type ErrDoesNotExist struct {
	Value string
}

func (e ErrDoesNotExist) Error() string {
	return fmt.Sprintf("kennung does not exist: %s", e.Value)
}

func (e ErrDoesNotExist) Is(target error) bool {
	_, ok := target.(ErrDoesNotExist)
	return ok
}

type ErrHinweisenExhausted struct{}

func (e ErrHinweisenExhausted) Error() string {
	return "hinweisen exhausted"
}

func (e ErrHinweisenExhausted) Is(target error) bool {
	_, ok := target.(ErrHinweisenExhausted)
	return ok
}
