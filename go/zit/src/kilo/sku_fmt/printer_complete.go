package sku_fmt

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
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
		pool: pool.MakePool[sku.Transacted](
			nil,
			nil,
		),
	}

	go func(s *WriterComplete) {
		for sk := range s.chTransacted {
			ui.TodoP4("handle write errors")
			s.wBuf.WriteString(sk.GetObjectId().String())
			s.wBuf.WriteByte('\t')

			g := sk.GetObjectId().GetGenre()
			s.wBuf.WriteString(g.String())

			if g == genres.Zettel {
				s.wBuf.WriteString(": ")
				s.wBuf.WriteString(sk.GetType().String())
				s.wBuf.WriteString(" ")
				s.wBuf.WriteString(sk.GetMetadata().Description.String())
			}

			s.wBuf.WriteString("\n")
			w1.pool.Put(sk)
		}

		s.chDone <- struct{}{}
	}(w1)

	return w1
}

func (w *WriterComplete) WriteOneSkuType(
	co sku.SkuType,
) (err error) {
	switch co.GetState() {
	// case checked_out_state.Internal:
	// 	sku.Resetter.ResetWith(sk, co.GetSku())

	default:
		// sku.Resetter.ResetWith(sk, co.GetSkuExternal())
		// TODO use proper states
		// sku.Resetter.ResetWith(sk, co.GetSku())
	}

	if err = w.WriteOneTransacted(co.GetSku()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *WriterComplete) WriteOneTransacted(
	src *sku.Transacted,
) (err error) {
	if src.GetObjectId().String() == "/" {
		err = errors.New("empty sku")
		return
	}

	dst := w.pool.Get()
	sku.Resetter.ResetWith(dst, src)

	select {
	case <-w.chDone:
		err = collections.MakeErrStopIteration()

	case w.chTransacted <- dst:
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
