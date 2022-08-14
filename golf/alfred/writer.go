package alfred

import (
	"io"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/alfred"
	"github.com/friedenberg/zit/delta/etikett"
	"github.com/friedenberg/zit/delta/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Writer interface {
	WriteZettel(stored_zettel.Named) (n int, err error)
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

func (w *writer) WriteZettel(z stored_zettel.Named) (n int, err error) {
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
