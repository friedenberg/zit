package zettel

import (
	"fmt"
)

type ErrHasInvalidAkteShaOrFilePath struct {
	Value string
}

func (e ErrHasInvalidAkteShaOrFilePath) Error() string {
	return fmt.Sprintf(
		"zettel text has invalid akte sha or file path: %q",
		e.Value,
	)
}

func (e ErrHasInvalidAkteShaOrFilePath) Is(target error) (ok bool) {
	_, ok = target.(ErrHasInvalidAkteShaOrFilePath)
	return
}
