package alfred

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/alfred"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/zettel_named"
)

type Writer interface {
	WriteZettel(zettel_named.Named) (n int, err error)
	WriteEtikett(e etikett.Etikett) (n int, err error)
	WriteHinweis(e hinweis.Hinweis) (n int, err error)
	WriteError(in error) (n int, out error)
	Close() error
}

type writer struct {
	alfredWriter alfred.Writer
}

func NewWriter(out io.Writer) (w *writer, err error) {
	var aw alfred.Writer

	if aw, err = alfred.NewWriter(out); err != nil {
		err = errors.Error(err)
		return
	}

	w = &writer{
		alfredWriter: aw,
	}

	return
}

func (w *writer) WriteZettel(z zettel_named.Named) (n int, err error) {
	item := ZettelToItem(z)
	return w.alfredWriter.WriteItem(item)
}

func (w *writer) WriteEtikett(e etikett.Etikett) (n int, err error) {
	item := EtikettToItem(e)
	return w.alfredWriter.WriteItem(item)
}

func (w *writer) WriteHinweis(e hinweis.Hinweis) (n int, err error) {
	item := HinweisToItem(e)
	return w.alfredWriter.WriteItem(item)
}

func (w *writer) WriteError(in error) (n int, out error) {
	if in == nil {
		return 0, nil
	}

	item := ErrorToItem(in)
	return w.alfredWriter.WriteItem(item)
}

func (w writer) Close() (err error) {
	return w.alfredWriter.Close()
}
