package zettel

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/golf/sku"
)

type WriterComplete struct {
	wBuf         *bufio.Writer
	chTransacted chan sku.TransactedZettel
	chDone       chan struct{}
}

func MakeWriterComplete(w io.Writer) WriterComplete {
	w1 := WriterComplete{
		chTransacted: make(chan sku.TransactedZettel),
		chDone:       make(chan struct{}),
		wBuf:         bufio.NewWriter(w),
	}

	go func(s *WriterComplete) {
		for z := range s.chTransacted {
			if z.GetKennung().String() == "/" {
				errors.Err().Printf("empty: %#v", z)
				continue
			}

			errors.TodoP4("handle write errors")
			s.wBuf.WriteString(z.GetKennung().String())
			s.wBuf.WriteString("\tZettel: !")
			s.wBuf.WriteString(z.GetTyp().String())
			s.wBuf.WriteString(" ")
			s.wBuf.WriteString(z.GetMetadatei().Bezeichnung.String())
			s.wBuf.WriteString("\n")
		}

		s.chDone <- struct{}{}
	}(&w1)

	return w1
}

func (w *WriterComplete) WriteZettelVerzeichnisse(
	z *sku.TransactedZettel,
) (err error) {
	select {
	case <-w.chDone:
		err = collections.MakeErrStopIteration()

	case w.chTransacted <- *z:
	}

	return
}

func (w *WriterComplete) Close() (err error) {
	close(w.chTransacted)
	<-w.chDone

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
