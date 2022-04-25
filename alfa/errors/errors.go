package errors

import (
	"log"

	"golang.org/x/xerrors"
)

func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

func Is(err, target error) bool {
	return xerrors.Is(err, target)
}

func PanicIfError(err interface{}) {
	if err == nil {
		return
	}

	switch t := err.(type) {
	case func() error:
		PanicIfError(t())
	case error:
		log.Output(2, t.Error())
		panic(t)
	}
}
