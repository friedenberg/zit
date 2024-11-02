package dir_layout

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type RelativePath interface {
	Rel(string) string
}

func (s DirLayout) Rel(
	p string,
) (out string) {
	out = p

	p1, _ := filepath.Rel(s.cwd, p)

	if p1 != "" {
		out = p1
	}

	return
}

func (s DirLayout) MakeRelativePathStringFormatWriter() interfaces.StringFormatWriter[string] {
	return relativePathStringFormatWriter(s)
}

type relativePathStringFormatWriter DirLayout

func (f relativePathStringFormatWriter) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	p string,
) (n int64, err error) {
	var n1 int

	{
		// if p, err = filepath.Rel(s.cwd, p); err != nil {
		// 	err = errors.Wrap(err)
		// 	return
		// }

		p1, _ := filepath.Rel(f.cwd, p)

		if p1 != "" {
			p = p1
		}
	}

	n1, err = w.WriteString(p)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
