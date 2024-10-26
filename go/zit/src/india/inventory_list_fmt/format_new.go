package inventory_list_fmt

import (
	"bufio"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
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
	w1 io.Writer,
) (n int64, err error) {
	bw := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, bw)

	defer o.Restore()

	var n1 int64
	var n2 int

	for {
		sk, ok := o.PopAndSave()

		if !ok {
			break
		}

		if sk.Metadata.Sha().IsNull() {
			err = errors.Errorf("empty sha: %s", sk)
			return
		}

		n1, err = s.Box.WriteStringFormat(bw, sk)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = fmt.Fprintf(bw, "\n")
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s VersionedFormatNew) WriteInventoryListObject(
	o *sku.Transacted,
	w1 io.Writer,
) (n int64, err error) {
	bw := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, bw)

	var n1 int64
	var n2 int

	n1, err = s.Box.WriteStringFormat(bw, o)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = fmt.Fprintf(bw, "\n")
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

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
	r1 io.Reader,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	r := bufio.NewReader(r1)

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

	return
}
