package objekte

import (
	"fmt"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/gattung"
)

type ErrUnsupportedFormatterValue interface {
	error
	GetFormatValue() string
	GetGattung() schnittstellen.GattungLike
}

func IsErrUnsupportedFormatterValue(err error) bool {
	var e ErrUnsupportedFormatterValue
	return errors.Is(err, e)
}

func MakeErrUnsupportedFormatterValue(
	formatValue string,
	g schnittstellen.GattungLike,
) error {
	return errors.Wrap(
		errUnsupportedFormatter{
			format:  formatValue,
			gattung: gattung.Must(g),
		},
	)
}

type errUnsupportedFormatter struct {
	format  string
	gattung gattung.Gattung
}

func (e errUnsupportedFormatter) Error() string {
	return fmt.Sprintf(
		"unsupported formatter value %q for gattung %s",
		e.format,
		e.gattung,
	)
}

func (e errUnsupportedFormatter) Is(err error) (ok bool) {
	_, ok = err.(errUnsupportedFormatter)
	return
}

func (e errUnsupportedFormatter) GetFormatValue() string {
	return e.format
}

func (e errUnsupportedFormatter) GetGattung() schnittstellen.GattungLike {
	return e.gattung
}

type ErrLockRequired struct {
	Operation string
}

func (e ErrLockRequired) Is(target error) bool {
	_, ok := target.(ErrLockRequired)
	return ok
}

func (e ErrLockRequired) Error() string {
	return fmt.Sprintf(
		"lock required for operation: %q",
		e.Operation,
	)
}

type VerlorenAndGefundenError interface {
	error
	AddToLostAndFound(string) (string, error)
}
