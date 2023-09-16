package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type ObjekteSaver interface {
	SaveObjekte(sku.SkuLikePtr) error
	SaveObjekteIncludeTai(sku.SkuLikePtr) error
}

type objekteSaver struct {
	writerFactory schnittstellen.ObjekteWriterFactory
	formatter     objekte_format.Format
	options       objekte_format.Options
}

func MakeObjekteSaver(
	// TODO-P1 add objekte index
	writerFactory schnittstellen.ObjekteWriterFactory,
	pmf objekte_format.Format,
	op objekte_format.Options,
) ObjekteSaver {
	if writerFactory == nil {
		panic("schnittstellen.ObjekteWriterFactory was nil")
	}

	if pmf == nil {
		panic("persisted_metadatei_format.Format was nil")
	}

	return objekteSaver{
		writerFactory: writerFactory,
		formatter:     pmf,
		options:       op,
	}
}

func (h objekteSaver) SaveObjekte(
	tl sku.SkuLikePtr,
) (err error) {
	var w sha.WriteCloser

	if w, err = h.writerFactory.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = h.formatter.FormatPersistentMetadatei(w, tl, h.options); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.Make(w.GetShaLike())

	log.Log().Printf(
		"saving objekte: %s -> %s",
		tl.GetKennungLike().GetGattung(),
		sh,
	)

	tl.SetObjekteSha(sh)

	return
}

func (h objekteSaver) SaveObjekteIncludeTai(
	tl sku.SkuLikePtr,
) (err error) {
	var w sha.WriteCloser

	if w, err = h.writerFactory.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = h.formatter.FormatPersistentMetadatei(
		w,
		tl,
		objekte_format.Options{IncludeTai: true},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.Make(w.GetShaLike())

	log.Log().Printf(
		"saving objekte with tai: %s -> %s",
		tl.GetKennungLike().GetGattung(),
		sh,
	)

	tl.SetObjekteSha(sh)

	return
}
