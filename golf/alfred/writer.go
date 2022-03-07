package alfred

import "io"

type Writer interface {
	WriteZettel(_NamedZettel) (n int, err error)
	WriteEtikett(e _Etikett) (n int, err error)
	WriteHinweis(e _Hinweis) (n int, err error)
	WriteError(in error) (n int, out error)
	Close() error
}

type writer struct {
	alfredWriter _AlfredWriter
}

func NewWriter(out io.Writer) (w *writer, err error) {
	var aw _AlfredWriter

	if aw, err = _AlfredNewWriter(out); err != nil {
		err = _Error(err)
		return
	}

	w = &writer{
		alfredWriter: aw,
	}

	return
}

func (w *writer) WriteZettel(z _NamedZettel) (n int, err error) {
	item := ZettelToItem(z)
	return w.alfredWriter.WriteItem(item)
}

func (w *writer) WriteEtikett(e _Etikett) (n int, err error) {
	item := EtikettToItem(e)
	return w.alfredWriter.WriteItem(item)
}

func (w *writer) WriteHinweis(e _Hinweis) (n int, err error) {
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
