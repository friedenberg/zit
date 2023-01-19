package alfred

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type Writer struct {
	alfredWriter *alfred.Writer
	Abbr         schnittstellen.FuncAbbreviateKorper
}

func New(out io.Writer, ha schnittstellen.FuncAbbreviateKorper) (w *Writer, err error) {
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

func (w *Writer) WriteHinweis(e hinweis.Hinweis) (n int, err error) {
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
