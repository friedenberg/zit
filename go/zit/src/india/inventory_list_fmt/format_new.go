package inventory_list_fmt

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

type VersionedFormatNew struct {
	Box *box_format.Box
}

func (v VersionedFormatNew) GetVersionedFormat() VersionedFormat {
	return v
}

func (v VersionedFormatNew) makePrinter(
	out interfaces.WriterAndStringWriter,
) interfaces.FuncIter[*sku.Transacted] {
	return string_format_writer.MakeDelim(
		"\n",
		out,
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return v.Box.WriteStringFormat(w, o)
			},
		),
	)
}

func (s VersionedFormatNew) WriteInventoryListBlob(
	o *sku.List,
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

func (s VersionedFormatNew) WriteInventoryListObject(
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

func (s VersionedFormatNew) ReadInventoryListObject(
	r1 io.Reader,
) (n int64, o *sku.Transacted, err error) {
	o = sku.GetTransactedPool().Get()

	r := bufio.NewReader(r1)

	if n, err = s.Box.ReadStringFormat(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s VersionedFormatNew) StreamInventoryListBlobSkus(
	rf func(interfaces.ShaGetter) (interfaces.ShaReadCloser, error),
	blobSha interfaces.Sha,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var ar interfaces.ShaReadCloser

	if ar, err = rf(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	r := bufio.NewReader(ar)

	for {
		o := sku.GetTransactedPool().Get()

		if _, err = s.Box.ReadStringFormat(r, o); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if err = f(o); err != nil {
			err = errors.Wrapf(err, "Object: %s", o)
			return
		}
	}

	sh := ar.GetShaLike()

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
