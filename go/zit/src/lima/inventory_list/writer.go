package inventory_list

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_fmt"
)

type VersionedFormat interface {
	WriteInventoryListBlob(*sku.List, func() (sha.WriteCloser, error)) (*sha.Sha, error)
	WriteInventoryListObject(*sku.Transacted, func() (sha.WriteCloser, error)) (*sha.Sha, error)
	ReadInventoryListObject(io.Reader) (int64, *sku.Transacted, error)
	StreamInventoryListBlobSkus(
		rf func(interfaces.ShaGetter) (interfaces.ShaReadCloser, error),
		blobSha interfaces.Sha,
		f interfaces.FuncIter[*sku.Transacted],
	) error
}

type versionedFormatOld struct {
	inventory_list_fmt.Factory
}

func (v versionedFormatOld) GetVersionedFormat() VersionedFormat {
	return v
}

func (s versionedFormatOld) WriteInventoryListObject(
	o *sku.Transacted,
	wf func() (sha.WriteCloser, error),
) (sh *sha.Sha, err error) {
	var w sha.WriteCloser

	if w, err = wf(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = s.Format.FormatPersistentMetadata(
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

func (s versionedFormatOld) WriteInventoryListBlob(
	o *sku.List,
	wf func() (sha.WriteCloser, error),
) (sh *sha.Sha, err error) {
	var sw sha.WriteCloser

	if sw, err = wf(); err != nil {
		err = errors.Wrap(err)
		return
	}

	func() {
		defer errors.DeferredCloser(&err, sw)

		fo := s.MakePrinter(sw)

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

func (s versionedFormatOld) ReadInventoryListObject(
	r io.Reader,
) (n int64, o *sku.Transacted, err error) {
	o = sku.GetTransactedPool().Get()

	if n, err = s.Format.ParsePersistentMetadata(
		catgut.MakeRingBuffer(r, 0),
		o,
		s.Options,
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

func (s versionedFormatOld) StreamInventoryListBlobSkus(
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

	dec := s.MakeScanner(ar)

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
