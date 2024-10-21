package inventory_list

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

type versionedFormatNew struct {
	box *box_format.Box
}

func (v versionedFormatNew) GetVersionedFormat() versionedFormat {
	return v
}

func (v versionedFormatNew) makePrinter(
	out interfaces.WriterAndStringWriter,
) interfaces.FuncIter[*sku.Transacted] {
	return string_format_writer.MakeDelim(
		"\n",
		out,
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return v.box.WriteStringFormat(w, o)
			},
		),
	)
}

func (s versionedFormatNew) writeInventoryListBlob(
	o *InventoryList,
	wf func() (sha.WriteCloser, error),
) (sh *sha.Sha, err error) {
	var sw sha.WriteCloser

	if sw, err = wf(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, sw)

	bw := bufio.NewWriter(sw)
	defer errors.DeferredFlusher(&err, bw)

	fo := s.makePrinter(bw)

	defer o.Restore()

	for {
		sk, ok := o.PopAndSave()

		if !ok {
			break
		}

		if sk.Metadata.Sha().IsNull() {
			err = errors.Errorf("empty sha: %s", sk)
			return
		}

		if err = fo(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	sh = sha.Make(sw.GetShaLike())

	return
}

func (s versionedFormatNew) writeInventoryListObject(
	o *sku.Transacted,
	wf func() (sha.WriteCloser, error),
) (sh *sha.Sha, err error) {
	var w sha.WriteCloser

	if w, err = wf(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	fo := s.makePrinter(bw)

	if err = fo(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(w.GetShaLike())

	return
}

func (s versionedFormatNew) readInventoryListObject(
	r1 io.Reader,
) (n int64, o *sku.Transacted, err error) {
	o = sku.GetTransactedPool().Get()

	r := catgut.MakeRingBuffer(r1, 0)

	if n, err = s.box.ReadStringFormat(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s versionedFormatNew) readInventoryListBlob(
	rf func(interfaces.ShaGetter) (interfaces.ShaReadCloser, error),
	blobSha interfaces.Sha,
	a *InventoryList,
) (err error) {
	var ar interfaces.ShaReadCloser

	if ar, err = rf(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	sw := sha.MakeWriter(nil)

	var eof bool

	r := catgut.MakeRingBuffer(ar, 0)

	for !eof {
		o := sku.GetTransactedPool().Get()

		if _, err = s.box.ReadStringFormat(r, o); err != nil {
			if errors.IsEOF(err) {
				// TODO determine if EOF is after content or during
				eof = true
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		a.Add(o)
	}

	sh := sw.GetShaLike()

	if !sh.EqualsSha(blobSha) {
		err = errors.Errorf(
			"objekte had blob sha %s while blob reader had %s",
			blobSha,
			sh,
		)
		return
	}

	return
}
