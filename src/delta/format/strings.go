package format

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

const (
	StringDRArrow        = "↳"
	StringNew            = "new"
	StringSame           = "same"
	StringChanged        = "changed"
	StringDeleted        = "deleted"
	StringUpdated        = "updated"
	StringArchived       = "archived"
	StringUnchanged      = "unchanged"
	StringRecognized     = "recognized"
	StringCheckedOut     = "checked out"
	StringWouldDelete    = "would delete"
	StringUnrecognized   = "unrecognized"
	StringFormatDateTime = "06-01-02 15:04:05"
	StringIndent         = "                 "
	LenStringMax         = len(StringIndent) // TODO-P4 use reflection?
)

func MakeFormatStringRightAligned(
	f string,
	args ...any,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (n int64, err error) {
		f = fmt.Sprintf(f+" ", args...)

		diff := LenStringMax + 1 - utf8.RuneCountInString(f)

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
