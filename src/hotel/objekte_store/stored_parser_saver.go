package objekte_store

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
)

type StoredParseSaver[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
] interface {
	ParseSaveStored(
		sem sku.ExternalMaybe[K, KPtr],
		t *objekte.External[O, OPtr, K, KPtr],
	) (err error)
}

type storedParserSaver[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
] struct {
	awf          schnittstellen.AkteWriterFactory
	akteParser   objekte.AkteParser[OPtr]
	objekteSaver ObjekteSaver2
}

func MakeStoredParseSaver[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
](
	owf schnittstellen.ObjekteIOFactory,
	awf schnittstellen.AkteIOFactory,
	akteParser objekte.AkteParser[OPtr],
	pmf persisted_metadatei_format.Format,
) storedParserSaver[O, OPtr, K, KPtr] {
	if akteParser == nil {
		akteParser = MakeNopAkteFormat[O, OPtr](awf)
	}

	if pmf == nil {
		panic("persisted_metadatei_format.Format was nil")
	}

	return storedParserSaver[O, OPtr, K, KPtr]{
		awf:        awf,
		akteParser: akteParser,
		objekteSaver: MakeObjekteSaver2(
			owf,
			pmf,
		),
	}
}

func (h storedParserSaver[O, OPtr, K, KPtr]) ParseSaveStored(
	sem sku.ExternalMaybe[K, KPtr],
	t *objekte.External[O, OPtr, K, KPtr],
) (err error) {
	var f *os.File

	errors.TodoP2("support akte")
	if f, err = files.OpenExclusiveReadOnly(sem.FDs.Objekte.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.FDs = sem.FDs
	t.Sku.Kennung = sem.Kennung

	r := sha.MakeReadCloser(f)

	defer errors.DeferredCloser(&err, r)

	if _, err = h.readAkte(r, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetAkteSha(sha.Make(r.Sha()))

	if err = h.objekteSaver.SaveObjekte(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h storedParserSaver[O, OPtr, K, KPtr]) readAkte(
	r sha.ReadCloser,
	o OPtr,
) (n int64, err error) {
	if n, err = h.akteParser.ParseAkte(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d akte bytes", n)

	return
}
