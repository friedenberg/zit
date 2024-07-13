package umwelt

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
)

type ErrUnsupportedFormatterValue interface {
	error
	GetFormatValue() string
	GetGattung() interfaces.GattungLike
}

func IsErrUnsupportedFormatterValue(err error) bool {
	var e ErrUnsupportedFormatterValue
	return errors.Is(err, e)
}

func MakeErrUnsupportedFormatterValue(
	formatValue string,
	g interfaces.GattungLike,
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

func (e errUnsupportedFormatter) GetGattung() interfaces.GattungLike {
	return e.gattung
}
