package inventory_list_fmt

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type VersionedFormatOld struct {
	Factory
}

func (v VersionedFormatOld) GetVersionedFormat() VersionedFormat {
	return v
}

func (s VersionedFormatOld) WriteInventoryListObject(
	o *sku.Transacted,
	w io.Writer,
) (n int64, err error) {
	if n, err = s.Format.FormatPersistentMetadata(
		w,
		o,
		object_inventory_format.Options{Tai: true},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s VersionedFormatOld) WriteInventoryListBlob(
	o *sku.List,
	w io.Writer,
) (n int64, err error) {
	var n1 int64

	fo := s.MakePrinter(w)

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

		n1, err = fo.Print(sk)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s VersionedFormatOld) ReadInventoryListObject(
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

func (s VersionedFormatOld) StreamInventoryListBlobSkus(
	r1 io.Reader,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	dec := s.MakeScanner(r1)

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
}
