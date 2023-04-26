package remote_transfers

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type common struct {
	*umwelt.Umwelt
}

func (c common) GetInheritorZettel(
	orf schnittstellen.ObjekteReaderFactory,
	arf schnittstellen.AkteReaderFactory,
) objekte_store.TransactedInheritor {
	p := collections.MakePool[zettel.Transacted, *zettel.Transacted]()

	inflator := objekte_store.MakeTransactedInflator[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	](
		schnittstellen.MakeBespokeObjekteReadWriterFactory(
			orf,
			c.StoreObjekten().Zettel(),
		),
		schnittstellen.MakeBespokeAkteReadWriterFactory(
			arf,
			c.StoreObjekten(),
		),
		persisted_metadatei_format.FormatForVersion(
			c.Konfig().GetStoreVersion(),
		),
		objekte_store.MakeAkteFormat[zettel.Objekte, *zettel.Objekte](
			objekte.MakeNopAkteParseSaver[zettel.Objekte, *zettel.Objekte](
				c.StoreObjekten(),
			),
			nil,
			c.StoreObjekten(),
		),
		p,
	)

	return objekte_store.MakeTransactedInheritor[
		zettel.Transacted,
		*zettel.Transacted,
	](
		inflator,
		c.StoreObjekten().Zettel(),
		p,
	)
}

func (c common) GetInheritorTyp(
	orf schnittstellen.ObjekteReaderFactory,
	arf schnittstellen.AkteReaderFactory,
) objekte_store.TransactedInheritor {
	errors.TodoP1("fix issues with typ skus being delivered with empty data")
	p := collections.MakePool[typ.Transacted, *typ.Transacted]()

	inflator := objekte_store.MakeTransactedInflator[
		typ.Akte,
		*typ.Akte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Akte],
		*objekte.NilVerzeichnisse[typ.Akte],
	](
		schnittstellen.MakeBespokeObjekteReadWriterFactory(
			orf,
			c.StoreObjekten().Typ(),
		),
		schnittstellen.MakeBespokeAkteReadWriterFactory(
			arf,
			c.StoreObjekten(),
		),
		persisted_metadatei_format.FormatForVersion(
			c.Konfig().GetStoreVersion(),
		),
		objekte_store.MakeAkteFormat[typ.Akte, *typ.Akte](
			objekte.MakeTextParserIgnoreTomlErrors[typ.Akte](c.StoreObjekten()),
			objekte.ParsedAkteTomlFormatter[typ.Akte]{},
			c.StoreObjekten(),
		),
		p,
	)

	return objekte_store.MakeTransactedInheritor[typ.Transacted, *typ.Transacted](
		inflator,
		c.StoreObjekten().Typ(),
		p,
	)
}
