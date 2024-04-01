package catgut

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

var (
	stringSpace   = MakeFromString(" ")
	stringNewline = MakeFromString("\n")
)

func WriteKeySpaceValueNewline(
	w io.Writer,
	key, value io.WriterTo,
) (n int64, err error) {
	var n1 int64

	n1, err = key.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = stringSpace.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = value.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = stringNewline.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
