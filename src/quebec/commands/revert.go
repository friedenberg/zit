package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/lima/bestandsaufnahme"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Revert struct {
	Last bool
}

func init() {
	registerCommandWithQuery(
		"revert",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Revert{}

			f.BoolVar(&c.Last, "last", false, "revert the last changes")

			return c
		},
	)
}

func (c Revert) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Revert) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		// gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

type revertMutterToKennungTuple struct {
	*kennung.Kennung2
	*sha.Sha
}

func (c Revert) RunWithQuery(u *umwelt.Umwelt, ms matcher.Query) (err error) {
	var mutterToKennung []revertMutterToKennungTuple

	switch {
	case c.Last:
		mutterToKennung, err = c.muttersFromLast(u)

	default:
		mutterToKennung, err = c.muttersFromQuery(u, ms)
	}

	if err != nil {
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	for _, rt := range mutterToKennung {
		if err = u.StoreObjekten().SetTransactedTo(rt.Kennung2, rt.Sha); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Revert) addOneMutter(
	z *sku.Transacted,
	mutters *[]revertMutterToKennungTuple,
) (err error) {
	mu := &z.Metadatei.Mutter

	if mu.IsNull() {
		// log.Err().Printf("%s has null mutter, cannot revert", z)
		return
	}

	var k kennung.Kennung2

	if err = k.SetWithKennung(z.GetKennung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh sha.Sha

	if err = sh.SetShaLike(mu); err != nil {
		err = errors.Wrap(err)
		return
	}

	rmtkt := revertMutterToKennungTuple{
		Kennung2: &k,
		Sha:      &sh,
	}

	*mutters = append(*mutters, rmtkt)

	return
}

func (c Revert) muttersFromQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (mutterToKennung []revertMutterToKennungTuple, err error) {
	mutterToKennung = make([]revertMutterToKennungTuple, 0)

	if err = u.StoreObjekten().QueryWithoutCwd(
		ms,
		func(z *sku.Transacted) (err error) {
			return c.addOneMutter(z, &mutterToKennung)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Revert) muttersFromLast(
	u *umwelt.Umwelt,
) (mutterToKennung []revertMutterToKennungTuple, err error) {
	s := u.StoreObjekten()

	var b *sku.Transacted

	if b, err = s.GetBestandsaufnahmeStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var a *bestandsaufnahme.Akte

	if a, err = s.GetBestandsaufnahmeStore().GetAkte(b.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	mutterToKennung = make([]revertMutterToKennungTuple, 0, a.Skus.Len())

	var formatGeneric objekte_format.FormatGeneric

	if formatGeneric, err = objekte_format.FormatForKeyError("Metadatei"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.Skus.EachPtr(
		func(sk *sku.Transacted) (err error) {
			var sh *sha.Sha

			if sh, err = objekte_format.GetShaForMetadatei(
				formatGeneric,
				sk.GetMetadatei(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if sk, err = u.StoreUtil().GetVerzeichnisse().ReadOneShas(sh); err != nil {
				err = errors.Wrap(err)
				return
			}

			return c.addOneMutter(sk, &mutterToKennung)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
