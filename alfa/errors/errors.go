package errors

import "golang.org/x/xerrors"

func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

func PanicIfError(err interface{}) {
	if err == nil {
		return
	}

	switch t := err.(type) {
	case func() error:
		PanicIfError(t())
	case error:
		panic(err)
	}
}
