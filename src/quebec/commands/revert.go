package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
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

func (c Revert) RunWithQuery(u *umwelt.Umwelt, ms matcher.Query) (err error) {
	var mutterToKennung map[string]string

	switch {
	case c.Last:
		mutterToKennung, err = c.mutterToKennungFromLast(u)

	default:
		mutterToKennung, err = c.mutterToKennungFromQuery(u, ms)
	}

	if err != nil {
		return
	}

	mutters := sku.MakeTransactedMutableSet()

	if err = u.StoreObjekten().ReadAll(
		gattungen.MakeSet(gattung.TrueGattung()...),
		func(z *sku.Transacted) (err error) {
			ms := z.Metadatei.Sha.String()

			kin, ok := mutterToKennung[ms]

			if !ok {
				return
			}

			if kin != z.Kennung.String() {
				return
			}

			mu := sku.GetTransactedPool().Get()

			if err = mu.SetFromTransacted(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return mutters.Add(mu)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if err = u.StoreObjekten().UpdateManyMetadatei(mutters); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Revert) addOneMutterToKennung(
	z *sku.Transacted,
	mutterToKennung map[string]string,
) (err error) {
	mu := &z.Metadatei.Mutter

	if mu.IsNull() {
		// log.Err().Printf("%s has null mutter, cannot revert", z)
		return
	}

	mutterToKennung[mu.String()] = z.Kennung.String()

	return
}

func (c Revert) mutterToKennungFromQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (mutterToKennung map[string]string, err error) {
	mutterToKennung = make(map[string]string)

	if err = u.StoreObjekten().QueryWithoutCwd(
		ms,
		func(z *sku.Transacted) (err error) {
			return c.addOneMutterToKennung(z, mutterToKennung)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Revert) mutterToKennungFromLast(
	u *umwelt.Umwelt,
) (mutterToKennung map[string]string, err error) {
	mutterToKennung = make(map[string]string)

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

	if err = a.Skus.EachPtr(
		func(sk *sku.Transacted) (err error) {
			return c.addOneMutterToKennung(sk, mutterToKennung)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
