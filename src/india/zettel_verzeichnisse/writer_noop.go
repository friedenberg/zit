package zettel_verzeichnisse

type WriterNoop struct{}

func (w WriterNoop) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	return
}
