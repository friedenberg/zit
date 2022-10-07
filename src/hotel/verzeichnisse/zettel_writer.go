package verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type Writer interface {
	WriteZettelVerzeichnisse(z *Zettel) (err error)
}

type WriterZettelTransacted struct {
	zettel_transacted.Writer
}

func (w WriterZettelTransacted) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	return w.WriteZettelTransacted(&z.Transacted)
}

type WriterFunc func(*Zettel) error

type writerFunc WriterFunc

func (w writerFunc) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	return WriterFunc(w)(z)
}

func MakeWriter(f WriterFunc) Writer {
	return writerFunc(f)
}

type writer struct {
	writers []Writer
	*ZettelPool
}

func (w writer) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	if w.ZettelPool != nil {
		defer w.Put(z)
	}

	for _, w := range w.writers {
		if err = w.WriteZettelVerzeichnisse(z); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func MakeWriterNamedFilters(fs ...zettel_named.NamedFilter) Writer {
	return MakeWriter(
		func(zt *Zettel) (err error) {
			for _, q := range fs {
				if !q.IncludeNamedZettel(zt.Transacted.Named) {
					err = io.EOF
					return
				}
			}

			return
		},
	)
}
