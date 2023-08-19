package alfred

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/hotel/transacted"
)

type Writer struct {
	alfredWriter *alfred.Writer
	kennungIndex kennung_index.Index
	typenIndex   kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]
	Abbr         func(kennung.Hinweis) (string, error)
}

func New(
	out io.Writer,
	kennungIndex kennung_index.Index,
	typenIndex kennung_index.KennungIndex[kennung.Typ, *kennung.Typ],
	ha func(kennung.Hinweis) (string, error),
) (w *Writer, err error) {
	var aw *alfred.Writer

	if aw, err = alfred.NewWriter(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	if kennungIndex == nil {
		panic("kennungIndex was nil")
	}

	if typenIndex == nil {
		panic("typenIndex was nil")
	}

	w = &Writer{
		Abbr:         ha,
		kennungIndex: kennungIndex,
		typenIndex:   typenIndex,
		alfredWriter: aw,
	}

	return
}

func (w *Writer) WriteZettelVerzeichnisse(z *transacted.Zettel) (err error) {
	item := w.zettelToItem(z, w.Abbr)
	w.alfredWriter.WriteItem(item)

	return
}

func (w *Writer) WriteEtikett(
	e kennung.IndexedLike[kennung.Etikett, *kennung.Etikett],
) (n int64, err error) {
	item := w.etikettToItem(e)
	w.alfredWriter.WriteItem(item)
	return
}

func (w *Writer) WriteHinweis(e kennung.Hinweis) (n int64, err error) {
	item := w.hinweisToItem(e)
	w.alfredWriter.WriteItem(item)
	return
}

func (w *Writer) WriteError(in error) (n int64, out error) {
	if in == nil {
		return 0, nil
	}

	item := w.errorToItem(in)
	w.alfredWriter.WriteItem(item)
	return
}

func (w Writer) Close() (err error) {
	return w.alfredWriter.Close()
}
