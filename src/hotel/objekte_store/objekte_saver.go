package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
)

type ObjekteSaver interface {
	SaveObjekte(objekte.StoredLikePtr) error
}

type objekteSaver struct {
	writerFactory schnittstellen.ObjekteWriterFactory
	formatter     persisted_metadatei_format.Format
}

func MakeObjekteSaver(
	writerFactory schnittstellen.ObjekteWriterFactory,
	pmf persisted_metadatei_format.Format,
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
	}
}

func (h objekteSaver) SaveObjekte(
	tl objekte.StoredLikePtr,
) (err error) {
	var w sha.WriteCloser

	if w, err = h.writerFactory.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = h.formatter.FormatPersistentMetadatei(w, tl); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.Make(w.Sha())

	tl.SetObjekteSha(sh)

	return
}