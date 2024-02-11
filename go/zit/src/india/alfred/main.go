package alfred

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/alfred"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

type Writer struct {
	alfredWriter *alfred.Writer
	Abbr         func(*kennung.Hinweis) (string, error)
}

func New(
	out io.Writer,
	ha func(*kennung.Hinweis) (string, error),
) (w *Writer, err error) {
	var aw *alfred.Writer

	if aw, err = alfred.NewWriter(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = &Writer{
		Abbr:         ha,
		alfredWriter: aw,
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
