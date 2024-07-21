package sku_fmt

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type WriterComplete struct {
	wBuf         *bufio.Writer
	pool         interfaces.Pool[sku.Transacted, *sku.Transacted]
	chTransacted chan *sku.Transacted
	chDone       chan struct{}
}

func MakeWriterComplete(w io.Writer) *WriterComplete {
	w1 := &WriterComplete{
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
			s.wBuf.WriteString(z.GetObjectId().String())
			s.wBuf.WriteByte('\t')

			g := z.GetObjectId().GetGenre()
			s.wBuf.WriteString(z.GetObjectId().GetGenre().String())

			if g == genres.Zettel {
				s.wBuf.WriteString(": !")
				s.wBuf.WriteString(z.GetType().String())
				s.wBuf.WriteString(" ")
				s.wBuf.WriteString(z.GetMetadata().Description.String())
			} else {
				s.wBuf.WriteString(g.String())
			}

			s.wBuf.WriteString("\n")
			w1.pool.Put(z)
		}

		s.chDone <- struct{}{}
	}(w1)

	return w1
}

func (w *WriterComplete) WriteOne(
	z *sku.Transacted,
) (err error) {
	if z.GetObjectId().String() == "/" {
		err = errors.New("empty sku")
		return
	}

	sk := w.pool.Get()
	sku.Resetter.ResetWith(sk, z)

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
