package stream_index

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
)

type PageId = sha.PageId

type Page struct {
	PageId
	*probe_index
	// All, Schwanzen  Page
	added, addedLatest *sku.TransactedHeap
	flushMode          object_mode.Mode
	hasChanges         bool
	directoryLayout    dir_layout.DirLayout
	config             *config.Compiled
	oids               map[string]struct{}
}

func (pt *Page) initialize(
	pid PageId,
	i *Index,
) {
	pt.directoryLayout = i.directoryLayout.SansObjectAge().SansObjectCompression()
	pt.PageId = pid
	pt.added = sku.MakeTransactedHeap()
	pt.addedLatest = sku.MakeTransactedHeap()
	pt.config = i.mutable_config
	pt.probe_index = &i.probe_index
	pt.oids = make(map[string]struct{})
}

func (s *Page) readOneRange(
	ra object_probe_index.Range,
	sk *sku.Transacted,
) (err error) {
	var f *os.File

	if f, err = files.Open(s.Path()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	b := make([]byte, ra.ContentLength)

	if _, err = f.ReadAt(b, ra.Offset); err != nil {
		err = errors.Wrap(err)
		return
	}

	dec := makeBinaryWithQueryGroup(nil, ids.SigilHistory)

	skWR := skuWithRangeAndSigil{
		skuWithSigil: skuWithSigil{
			Transacted: sk,
		},
		Range: ra,
	}

	if _, err = dec.readFormatExactly(f, &skWR); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (pt *Page) add(
	z1 *sku.Transacted,
	options sku.CommitOptions,
) (err error) {
	pt.oids[z1.ObjectId.String()] = struct{}{}
	z := sku.GetTransactedPool().Get()

	sku.TransactedResetter.ResetWith(z, z1)

	if options.Contains(object_mode.ModeLatest) && !options.ChangeIsHistorical {
		pt.addedLatest.Add(z)
	} else {
		pt.added.Add(z)
	}

	pt.hasChanges = true

	return
}

func (pt *Page) waitingToAddLen() int {
	return pt.added.Len() + pt.addedLatest.Len()
}

func (pt *Page) CopyJustHistory(
	s sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return pt.copyHistoryAndMaybeLatest(s, w, false, false)
}

func (pt *Page) copyJustHistoryFrom(
	r io.Reader,
	s sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[skuWithRangeAndSigil],
) (err error) {
	dec := makeBinaryWithQueryGroup(s, ids.SigilHistory)

	var sk skuWithRangeAndSigil

	for {
		sk.Offset += sk.ContentLength
		sk.Transacted = sku.GetTransactedPool().Get()
		sk.ContentLength, err = dec.readFormatAndMatchSigil(r, &sk)
		if err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if err = w(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}
}

func (pt *Page) copyJustHistoryAndAdded(
	s sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return pt.copyHistoryAndMaybeLatest(s, w, true, false)
}

func (pt *Page) copyHistoryAndMaybeLatest(
	qg sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
	includeAdded bool,
	includeAddedLatest bool,
) (err error) {
	var r io.ReadCloser

	if r, err = pt.directoryLayout.ReadCloserCache(pt.Path()); err != nil {
		if errors.IsNotExist(err) {
			r = io.NopCloser(bytes.NewReader(nil))
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.DeferredCloser(&err, r)

	br := bufio.NewReader(r)

	if !includeAdded && !includeAddedLatest {
		if err = pt.copyJustHistoryFrom(
			br,
			qg,
			func(sk skuWithRangeAndSigil) (err error) {
				if err = w(sk.Transacted); err != nil {
					err = errors.Wrapf(err, "%s", sk.Transacted)
					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	dec := makeBinaryWithQueryGroup(qg, ids.SigilHistory)

	ui.TodoP3("determine performance of this")
	added := pt.added.Copy()

	var sk skuWithRangeAndSigil

	if err = heap.MergeStream(
		&added,
		func() (tz *sku.Transacted, err error) {
			tz = sku.GetTransactedPool().Get()
			sk.Transacted = tz
			_, err = dec.readFormatAndMatchSigil(br, &sk)
			if err != nil {
				if errors.IsEOF(err) {
					err = collections.MakeErrStopIteration()
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			return
		},
		w,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !includeAddedLatest {
		return
	}

	addedSchwanz := pt.addedLatest.Copy()

	if err = heap.MergeStream(
		&addedSchwanz,
		func() (tz *sku.Transacted, err error) {
			err = collections.MakeErrStopIteration()
			return
		},
		w,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (pt *Page) MakeFlush(
	changesAreHistorical bool,
) func() error {
	return func() (err error) {
		pw := &writer{
			Page:        pt,
			probe_index: pt.probe_index,
		}

		if changesAreHistorical {
			pw.changesAreHistorical = true
			pw.hasChanges = true
		}

		if err = pw.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}

		pt.hasChanges = false

		return
	}
}
