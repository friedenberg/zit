package errors

import "golang.org/x/xerrors"

type (
	_Frame   = xerrors.Frame
	_Printer = xerrors.Printer
)

var (
	_Caller      = xerrors.Caller
	_FormatError = xerrors.FormatError
	_Errorf      = xerrors.Errorf
)
