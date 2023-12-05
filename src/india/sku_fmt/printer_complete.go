package sku_fmt

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type WriterComplete struct {
	wBuf         *bufio.Writer
	pool         schnittstellen.Pool[sku.Transacted, *sku.Transacted]
	chTransacted chan *sku.Transacted
	chDone       chan struct{}
}

func MakeWriterComplete(w io.Writer) WriterComplete {
	w1 := WriterComplete{
		chTransacted: make(chan *sku.Transacted),
		chDone:       make(chan struct{}),
		wBuf:         bufio.NewWriter(w),
		pool: pool.MakePool[sku.Transacted, *sku.Transacted](
			nil,
			nil,
		),
	}

	go func(s *WriterComplete) {
		for z := range s.chTransacted {
			errors.TodoP4("handle write errors")
			s.wBuf.WriteString(z.GetKennung().String())
			s.wBuf.WriteString("\tZettel: !")
			s.wBuf.WriteString(z.GetTyp().String())
			s.wBuf.WriteString(" ")
			s.wBuf.WriteString(z.GetMetadatei().Bezeichnung.String())
			s.wBuf.WriteString("\n")
			w1.pool.Put(z)
		}

		s.chDone <- struct{}{}
	}(&w1)

	return w1
}

func (w *WriterComplete) WriteZettelVerzeichnisse(
	z *sku.Transacted,
) (err error) {
	if z.GetKennung().String() == "/" {
		err = errors.New("empty sku")
		return
	}

	sk := w.pool.Get()

	if err = sk.SetFromSkuLike(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	select {
	case <-w.chDone:
		err = collections.MakeErrStopIteration()

	case w.chTransacted <- sk:
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
