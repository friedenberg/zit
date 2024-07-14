package env

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

type ErrUnsupportedFormatterValue interface {
	error
	GetFormatValue() string
	interfaces.GenreGetter
}

func IsErrUnsupportedFormatterValue(err error) bool {
	var e ErrUnsupportedFormatterValue
	return errors.Is(err, e)
}

func MakeErrUnsupportedFormatterValue(
	formatValue string,
	g interfaces.Genre,
) error {
	return errors.Wrap(
		errUnsupportedFormatter{
			format: formatValue,
			genres: genres.Must(g),
		},
	)
}

type errUnsupportedFormatter struct {
	format string
	genres genres.Genre
}

func (e errUnsupportedFormatter) Error() string {
	return fmt.Sprintf(
		"unsupported formatter value %q for gattung %s",
		e.format,
		e.genres,
	)
}

func (e errUnsupportedFormatter) Is(err error) (ok bool) {
	_, ok = err.(errUnsupportedFormatter)
	return
}

func (e errUnsupportedFormatter) GetFormatValue() string {
	return e.format
}

func (e errUnsupportedFormatter) GetGattung() interfaces.Genre {
	return e.genres
}
