package standort

import (
	"io"
	"path/filepath"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

func (s Standort) Rel(
	p string,
) (out string) {
	out = p

	p1, _ := filepath.Rel(s.cwd, p)

	if p1 != "" {
		out = p1
	}

	return
}

func (s Standort) MakeRelativePathStringFormatWriter() schnittstellen.StringFormatWriter[string] {
	return relativePathStringFormatWriter(s)
}

type relativePathStringFormatWriter Standort

func (f relativePathStringFormatWriter) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
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

func (s Standort) MakeWriterRelativePath(
	p string,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (n int64, err error) {
		var n1 int

		{
			// if p, err = filepath.Rel(s.cwd, p); err != nil {
			// 	err = errors.Wrap(err)
			// 	return
			// }

			p1, _ := filepath.Rel(s.cwd, p)

			if p1 != "" {
				p = p1
			}
		}

		if n1, err = io.WriteString(w, p); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}

func (s Standort) MakeWriterRelativePathOr(
	p string,
	or schnittstellen.FuncWriter,
) schnittstellen.FuncWriter {
	if p == "" {
		return or
	}

	return func(w io.Writer) (n int64, err error) {
		var n1 int

		{
			// if p, err = filepath.Rel(s.cwd, p); err != nil {
			// 	err = errors.Wrap(err)
			// 	return
			// }

			p1, _ := filepath.Rel(s.cwd, p)

			if p1 != "" {
				p = p1
			}
		}

		if n1, err = io.WriteString(w, p); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}
