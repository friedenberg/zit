package alfred

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/alfred"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt_debug"
)

type Writer struct {
	alfredWriter *alfred.Writer
	abbr         ids.Abbr
	organizeFmt  interfaces.StringFormatWriter[sku.ExternalLike]
}

func New(
	out io.Writer,
	abbr ids.Abbr,
	organizeFmt interfaces.StringFormatWriter[sku.ExternalLike],
) (w *Writer, err error) {
	var aw *alfred.Writer

	if aw, err = alfred.NewWriter(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = &Writer{
		abbr:         abbr,
		alfredWriter: aw,
		organizeFmt:  organizeFmt,
	}

	return
}

func (w *Writer) PrintOne(z *sku.Transacted) (err error) {
	var item *alfred.Item
	g := z.GetGenre()

	switch g {
	case genres.Zettel:
		item = w.zettelToItem(z)

	case genres.Tag:
		var e ids.Tag

		if err = e.Set(z.ObjectId.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		item = w.etikettToItem(z, &e)

	default:
		item = w.alfredWriter.Get()
		item.Title = fmt.Sprintf("not implemented for gattung: %q", g)
		item.Subtitle = fmt.Sprintf("%s", sku_fmt_debug.StringTaiGenreObjectIdShaBlob(z))
	}

	w.alfredWriter.WriteItem(item)

	return
}

func (w *Writer) WriteHinweis(e ids.ZettelId) (n int64, err error) {
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
