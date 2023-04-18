package remote_transfers

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte"
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
			c.StoreObjekten(),
		),
		schnittstellen.MakeBespokeAkteReadWriterFactory(
			arf,
			c.StoreObjekten(),
		),
		metadatei.PersistedFormat{},
		objekte.MakeNopAkteFormat[zettel.Objekte, *zettel.Objekte](),
		p,
	)

	return objekte_store.MakeTransactedInheritor[zettel.Transacted, *zettel.Transacted](
		inflator,
		c.StoreObjekten().Zettel(),
		p,
	)
}

func (c common) GetInheritorTyp(
	orf schnittstellen.ObjekteReaderFactory,
	arf schnittstellen.AkteReaderFactory,
) objekte_store.TransactedInheritor {
	errors.TodoP0("fix issues with typ skus being delivered with empty data")
	p := collections.MakePool[typ.Transacted, *typ.Transacted]()

	inflator := objekte_store.MakeTransactedInflator[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	](
		schnittstellen.MakeBespokeObjekteReadWriterFactory(
			orf,
			c.StoreObjekten(),
		),
		schnittstellen.MakeBespokeAkteReadWriterFactory(
			arf,
			c.StoreObjekten(),
		),
		metadatei.PersistedFormat{},
		typ.MakeFormatTextIgnoreTomlErrors(c.StoreObjekten()),
		p,
	)

	return objekte_store.MakeTransactedInheritor[typ.Transacted, *typ.Transacted](
		inflator,
		c.StoreObjekten().Typ(),
		p,
	)
}
