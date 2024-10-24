package inventory_list

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_fax"
)

type versionedFormat interface {
	writeInventoryListBlob(*InventoryList, func() (sha.WriteCloser, error)) (*sha.Sha, error)
	writeInventoryListObject(*sku.Transacted, func() (sha.WriteCloser, error)) (*sha.Sha, error)
	readInventoryListObject(io.Reader) (int64, *sku.Transacted, error)
	streamInventoryListBlobSkus(
		rf func(interfaces.ShaGetter) (interfaces.ShaReadCloser, error),
		blobSha interfaces.Sha,
		f interfaces.FuncIter[*sku.Transacted],
	) error
}

type versionedFormatOld struct {
	object_format object_inventory_format.Format
	options       object_inventory_format.Options
	format
}

func (v versionedFormatOld) GetVersionedFormat() versionedFormat {
	return v
}

func (s versionedFormatOld) writeInventoryListObject(
	o *sku.Transacted,
	wf func() (sha.WriteCloser, error),
) (sh *sha.Sha, err error) {
	var w sha.WriteCloser

	if w, err = wf(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = s.object_format.FormatPersistentMetadata(
		w,
		o,
		object_inventory_format.Options{Tai: true},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(w.GetShaLike())

	return
}

func (s versionedFormatOld) writeInventoryListBlob(
	o *InventoryList,
	wf func() (sha.WriteCloser, error),
) (sh *sha.Sha, err error) {
	var sw sha.WriteCloser

	if sw, err = wf(); err != nil {
		err = errors.Wrap(err)
		return
	}

	func() {
		defer errors.DeferredCloser(&err, sw)

		fo := inventory_list_fax.MakePrinter(
			sw,
			s.object_format,
			s.options,
		)

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

			_, err = fo.Print(sk)
			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}()

	sh = sha.Make(sw.GetShaLike())

	return
}

func (s versionedFormatOld) readInventoryListObject(
	r io.Reader,
) (n int64, o *sku.Transacted, err error) {
	o = sku.GetTransactedPool().Get()

	if n, err = s.object_format.ParsePersistentMetadata(
		catgut.MakeRingBuffer(r, 0),
		o,
		s.options,
	); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s versionedFormatOld) streamInventoryListBlobSkus(
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

	sw := sha.MakeWriter(nil)

	dec := inventory_list_fax.MakeScanner(
		ar,
		s.format,
		s.options,
	)

	// dec.SetDebug()

	for dec.Scan() {
		sk := dec.GetTransacted()

		if err = f(sk); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}
	}

	if err = dec.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return

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
