package objekte

import "github.com/friedenberg/zit/src/bravo/line_format"

type WriterLineFormat struct {
	lf *line_format.Writer
}

func MakeWriterLineFormat(lf *line_format.Writer) WriterLineFormat {
	return WriterLineFormat{
		lf: lf,
	}
}

func (w WriterLineFormat) WriteObjekte(o Objekte) (err error) {
	w.lf.WriteFormat(
		"%s %s %s %s %s",
		o.Gattung,
		o.Mutter[0],
		o.Mutter[1],
		o.Id,
		o.Sha,
	)

	return
}
