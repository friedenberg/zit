package objekte

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
)

type ErrUnsupportedFormatterValue interface {
	error
	GetFormatValue() string
	GetGattung() schnittstellen.Gattung
}

func IsErrUnsupportedFormatterValue(err error) bool {
	var e ErrUnsupportedFormatterValue
	return errors.Is(err, e)
}

func MakeErrUnsupportedFormatterValue(
	formatValue string,
	g schnittstellen.Gattung,
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

func (e errUnsupportedFormatter) GetGattung() schnittstellen.Gattung {
	return e.gattung
}