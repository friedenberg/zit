package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
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

type revertTuple struct {
	*sku.Transacted
	*sha.Sha
}

func (c Revert) RunWithQuery(u *umwelt.Umwelt, ms matcher.Query) (err error) {
	f := func(rt revertTuple) (err error) {
		if rt.IsNull() {
			return
		}

		if err = u.StoreObjekten().RevertTo(rt.Transacted, rt.Sha); err != nil {
			err = errors.Wrapf(err, "Sha %s", rt.Sha)
			return
		}

		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	switch {
	case c.Last:
		err = c.runRevertFromLast(u, f)

	default:
		err = c.runRevertFromQuery(u, ms, f)
	}

	if err != nil {
		return
	}

	return
}

func (c Revert) runRevertFromQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
	f schnittstellen.FuncIter[revertTuple],
) (err error) {
	if err = u.StoreObjekten().QueryWithoutCwd(
		ms,
		func(z *sku.Transacted) (err error) {
			return f(revertTuple{Transacted: z, Sha: z.Metadatei.Mutter()})
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Revert) runRevertFromLast(
	u *umwelt.Umwelt,
	f schnittstellen.FuncIter[revertTuple],
) (err error) {
	s := u.StoreObjekten()

	var b *sku.Transacted

	if b, err = s.GetBestandsaufnahmeStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetBestandsaufnahmeStore().StreamAkte(
		b.GetAkteSha(),
		func(sk *sku.Transacted) (err error) {
			return f(revertTuple{Transacted: sk, Sha: sk.Metadatei.Mutter()})
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
