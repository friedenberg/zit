package alfred

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

type Writer struct {
	alfredWriter   *alfred.Writer
	kennungIndex   kennung_index.Index
	typenIndex     kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]
	etikettenIndex kennung_index.EtikettIndex
	Abbr           func(kennung.Hinweis) (string, error)
}

func New(
	out io.Writer,
	kennungIndex kennung_index.Index,
	typenIndex kennung_index.KennungIndex[kennung.Typ, *kennung.Typ],
	etikettenIndex kennung_index.EtikettIndex,
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
		Abbr:           ha,
		kennungIndex:   kennungIndex,
		typenIndex:     typenIndex,
		etikettenIndex: etikettenIndex,
		alfredWriter:   aw,
	}

	return
}

func (w *Writer) PrintOne(z *sku.Transacted) (err error) {
	var item *alfred.Item
	g := z.GetGattung()

	switch g {
	case gattung.Zettel:
		item = w.zettelToItem(z, w.Abbr)

	case gattung.Etikett:
		var e kennung.Etikett

		if err = e.Set(z.Kennung.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		item = w.etikettToItem(&e)

	default:
		item = w.alfredWriter.Get()
		item.Title = fmt.Sprintf("not implemented for gattung: %q", g)
		item.Subtitle = fmt.Sprintf("%s", sku_fmt.String(z))
	}

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