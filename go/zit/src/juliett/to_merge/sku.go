package to_merge

import (
	"bufio"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

type Sku struct {
	ConflictMarkerPath  string
	Left, Middle, Right *sku.Transacted
}

func (tm Sku) IsAllInlineTyp(itc kennung.InlineTypChecker) bool {
	if !itc.IsInlineTyp(tm.Left.GetTyp()) {
		return false
	}

	if !itc.IsInlineTyp(tm.Middle.GetTyp()) {
		return false
	}

	if !itc.IsInlineTyp(tm.Right.GetTyp()) {
		return false
	}

	return true
}

func (tm *Sku) MergeEtiketten() (err error) {
	left := tm.Left.GetEtiketten().CloneMutableSetPtrLike()
	middle := tm.Middle.GetEtiketten().CloneMutableSetPtrLike()
	right := tm.Right.GetEtiketten().CloneMutableSetPtrLike()

	same := kennung.MakeEtikettMutableSet()
	deleted := kennung.MakeEtikettMutableSet()

	removeFromAllButAddTo := func(
		e *kennung.Etikett,
		toAdd kennung.EtikettMutableSet,
	) (err error) {
		if err = toAdd.AddPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = left.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = middle.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = right.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = middle.EachPtr(
		func(e *kennung.Etikett) (err error) {
			if left.ContainsKey(left.KeyPtr(e)) && right.ContainsKey(right.KeyPtr(e)) {
				return removeFromAllButAddTo(e, same)
			} else if left.ContainsKey(left.KeyPtr(e)) || right.ContainsKey(right.KeyPtr(e)) {
				return removeFromAllButAddTo(e, deleted)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = left.EachPtr(same.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = right.EachPtr(same.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	ets := same.CloneSetPtrLike()

	tm.Left.GetMetadatei().SetEtiketten(ets)
	tm.Middle.GetMetadatei().SetEtiketten(ets)
	tm.Right.GetMetadatei().SetEtiketten(ets)

	return
}

func (tm *Sku) ReadConflictMarker(
	sv schnittstellen.StoreVersion,
	op objekte_format.Options,
) (err error) {
	var f *os.File

	if f, err = files.Open(tm.ConflictMarkerPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	br := bufio.NewReader(f)

	s := sku_fmt.MakeFormatBestandsaufnahmeScanner(
		br,
		objekte_format.FormatForVersion(sv),
		op,
	)

	i := 0

	for s.Scan() {
		sk := s.GetTransacted()

		switch i {
		case 0:
			tm.Left = sk

		case 1:
			tm.Middle = sk

		case 2:
			tm.Right = sk

		default:
			err = errors.Errorf("too many skus in conflict file")
			return
		}

		i++
	}

	if err = s.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (tm Sku) WriteConflictMarker(
	s standort.Standort,
	sv schnittstellen.StoreVersion,
	op objekte_format.Options,
	path string,
) (err error) {
	var f *os.File

	if f, err = s.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	bw := bufio.NewWriter(f)
	defer errors.DeferredFlusher(&err, bw)

	p := sku_fmt.MakeFormatBestandsaufnahmePrinter(
		bw,
		objekte_format.FormatForVersion(sv),
		op,
	)

	if _, err = p.PrintMany(tm.Left, tm.Middle, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Rename(f.Name(), path); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
