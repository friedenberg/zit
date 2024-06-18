package typ_akte

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type formatterFormatterUTIGroups struct{}

func MakeFormatterFormatterUTIGroups() *formatterFormatterUTIGroups {
	return &formatterFormatterUTIGroups{}
}

func (f formatterFormatterUTIGroups) Format(
	w io.Writer,
	ct *V0,
) (n int64, err error) {
	for groupName, group := range ct.FormatterUTIGroups {
		sb := bytes.NewBuffer(nil)

		sb.WriteString(groupName)

		for uti, formatter := range group.Map() {
			sb.WriteString(" ")
			sb.WriteString(uti)
			sb.WriteString(" ")
			sb.WriteString(formatter)
		}

		sb.WriteString("\n")

		if n, err = io.Copy(w, sb); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
