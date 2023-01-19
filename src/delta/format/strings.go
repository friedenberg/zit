package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

const (
	StringNew          = "new"
	StringSame         = "same"
	StringChanged      = "changed"
	StringDeleted      = "deleted"
	StringUpdated      = "updated"
	StringArchived     = "archived"
	StringUnchanged    = "unchanged"
	StringRecognized   = "recognized"
	StringCheckedOut   = "checked out"
	StringUnrecognized = "unrecognized"
	// StringHeaderIndent = "=============== "
	StringHeaderIndent = "                "
	//TODO-P4 use reflection?
	LenStringMax = len(StringUnrecognized) + 4
)

func MakeFormatStringRightAlignedParen(
	f string,
) schnittstellen.FuncWriter {
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

func MakeWriterFormatStringIndentedHeader(
	cw FuncColorWriter,
	indentString string,
) schnittstellen.FuncWriterFormat[string] {
	return func(w io.Writer, v string) (n int64, err error) {
		return Write(
			w,
			MakeFormatString(indentString),
			cw(MakeFormatString("%s", v), ColorTypeTitle),
		)
	}
}
