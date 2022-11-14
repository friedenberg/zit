package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

const (
	StringNew          = "new"
	StringSame         = "same"
	StringChanged      = "changed"
	StringDeleted      = "deleted"
	StringUpdated      = "updated"
	StringArchived     = "archived"
	StringRecognized   = "recognized"
	StringCheckedOut   = "checked out"
	StringUnrecognized = "unrecognized"
	LenStringMax       = len(StringUnrecognized) + 4
)

func MakeFormatStringRightAlignedParen(
	f string,
) WriterFunc {
	return func(w io.Writer) (n int64, err error) {
		f = fmt.Sprintf("(%s) ", f)

		diff := LenStringMax - len(f)

		if diff > 0 {
			f = strings.Repeat(" ", diff) + f
		}

		var n1 int

		if n1, err = io.WriteString(w, f); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}
