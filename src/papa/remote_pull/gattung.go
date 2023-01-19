package remote_pull

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

func (c *client) PullSkus(
	filter id_set.Filter,
	gattungSet gattungen.Set,
) (err error) {
	errors.TodoP0("implement etikett and akte")
	gattungInheritors := map[gattung.Gattung]objekte.TransactedInheritor{
		gattung.Zettel: c.GetInheritorZettel(),
		gattung.Typ:    c.GetInheritorTyp(),
	}

	if err = c.SkusFromFilter(
		filter,
		gattungSet,
		func(sk sku.Sku2) (err error) {
			var el objekte.TransactedInheritor
			ok := false

			if el, ok = gattungInheritors[sk.Gattung]; !ok {
				return
			}

			errors.TodoP2("check for akte sha")

			if c.umwelt.Standort().HasObjekte(sk.Gattung, sk.ObjekteSha) {
				errors.Log().Printf("already have objekte: %s", sk.ObjekteSha)
				return
			}

			errors.Log().Printf("need objekte: %s", sk.ObjekteSha)

			if err = el.InflateFromDataIdentityAndStoreAndInherit(sk); err != nil {
				err = errors.Wrapf(err, "Sku: %s", sk)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *client) GetInheritorZettel() objekte.TransactedInheritor {
	p := collections.MakePool2[zettel.Transacted, *zettel.Transacted]()

	inflator := objekte.MakeTransactedInflator[
		zettel.Objekte,
		*zettel.Objekte,
		hinweis.Hinweis,
		*hinweis.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	](
		schnittstellen.MakeBespokeObjekteReadWriterFactory(
			c,
			c.umwelt.StoreObjekten(),
		),
		schnittstellen.MakeBespokeAkteReadWriterFactory(
			c,
			c.umwelt.StoreObjekten(),
		),
		&zettel.FormatObjekte{
			IgnoreTypErrors: true,
		},
		objekte.MakeNopAkteFormat[zettel.Objekte, *zettel.Objekte](),
		p,
	)

	return objekte.MakeTransactedInheritor[zettel.Transacted, *zettel.Transacted](
		inflator,
		c.umwelt.StoreObjekten().Zettel(),
		p,
	)
}

func (c *client) GetInheritorTyp() objekte.TransactedInheritor {
	errors.TodoP0("fix issues with typ skus being delivered with empty data")
	p := collections.MakePool2[typ.Transacted, *typ.Transacted]()

	inflator := objekte.MakeTransactedInflator[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	](
		schnittstellen.MakeBespokeObjekteReadWriterFactory(
			c,
			c.umwelt.StoreObjekten(),
		),
		schnittstellen.MakeBespokeAkteReadWriterFactory(
			c,
			c.umwelt.StoreObjekten(),
		),
		nil,
		typ.MakeFormatTextIgnoreTomlErrors(c.umwelt.StoreObjekten()),
		p,
	)

	return objekte.MakeTransactedInheritor[typ.Transacted, *typ.Transacted](
		inflator,
		c.umwelt.StoreObjekten().Typ(),
		p,
	)
}
