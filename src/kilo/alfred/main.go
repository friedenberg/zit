package alfred

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type Writer struct {
	alfredWriter   *alfred.Writer
	etikettenIndex kennung_index.Index2[kennung.Etikett]
	typenIndex     kennung_index.Index2[kennung.Typ]
	Abbr           func(kennung.Hinweis) (string, error)
}

func New(
	out io.Writer,
	etikettenIndex kennung_index.Index2[kennung.Etikett],
	typenIndex kennung_index.Index2[kennung.Typ],
	ha func(kennung.Hinweis) (string, error),
) (w *Writer, err error) {
	var aw *alfred.Writer

	if aw, err = alfred.NewWriter(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	if etikettenIndex == nil {
		panic("etikettenIndex was nil")
	}

	if typenIndex == nil {
		panic("typenIndex was nil")
	}

	w = &Writer{
		Abbr:           ha,
		etikettenIndex: etikettenIndex,
		typenIndex:     typenIndex,
		alfredWriter:   aw,
	}

	return
}

func (w *Writer) WriteZettelVerzeichnisse(z *zettel.Transacted) (err error) {
	item := w.zettelToItem(z, w.Abbr)
	w.alfredWriter.WriteItem(item)

	return
}

func (w *Writer) WriteEtikett(e kennung.Etikett) (n int, err error) {
	item := w.etikettToItem(e)
	w.alfredWriter.WriteItem(item)
	return
}

func (w *Writer) WriteHinweis(e kennung.Hinweis) (n int, err error) {
	item := w.hinweisToItem(e)
	w.alfredWriter.WriteItem(item)
	return
}

func (w *Writer) WriteError(in error) (n int, out error) {
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
