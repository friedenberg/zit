package commands

import (
	"flag"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type CheckinAkte struct {
	Delete       bool
	NewEtiketten collections_ptr.Flag[kennung.Etikett, *kennung.Etikett]
}

func init() {
	registerCommand(
		"checkin-akte",
		func(f *flag.FlagSet) Command {
			c := &CheckinAkte{
				NewEtiketten: collections_ptr.MakeFlagCommas[kennung.Etikett](
					collections_ptr.SetterPolicyAppend,
				),
			}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.Var(
				c.NewEtiketten,
				"new-etiketten",
				"comma-separated etiketten (will replace existing Etiketten)",
			)

			return c
		},
	)
}

func (c CheckinAkte) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args)%2 != 0 {
		err = errors.Errorf(
			"arguments must come in pairs of hinweis and akte path",
		)
		return
	}

	type externalAktePair struct {
		*kennung.Hinweis
		path string
	}

	pairs := make([]externalAktePair, len(args)/2)

	// transform args into pairs of hinweis and filepaths
	for i, p := range pairs {
		hs := args[i*2]
		ap := args[(i*2)+1]

		if p.Hinweis, err = kennung.MakeHinweis(hs); err != nil {
			err = errors.Wrap(err)
			return
		}

		p.path = ap
		pairs[i] = p
	}

	zettels := make([]*sku.Transacted, len(pairs))

	// iterate through pairs and read current zettel
	for i, p := range pairs {
		if zettels[i], err = u.StoreObjekten().ReadOne(
			p.Hinweis,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for i, p := range pairs {
		var ow sha.WriteCloser

		if ow, err = u.Standort().AkteWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}

		var as sha.Sha

		shaError := as.Set(p.path)

		switch {
		case files.Exists(p.path):
			var f *os.File

			if f, err = files.Open(p.path); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, f)

			if _, err = io.Copy(ow, f); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = ow.Close(); err != nil {
				err = errors.Wrap(err)
				return
			}

			if zettels[i], err = u.StoreObjekten().ReadOne(
				p.Hinweis,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			zettels[i].SetAkteSha(ow.GetShaLike())

		case shaError == nil:
			zettels[i].SetAkteSha(&as)

		default:
			err = errors.Errorf("argument is neither sha nor path")
			return
		}

		if c.NewEtiketten.Len() > 0 {
			m := zettels[i].GetMetadatei()
			m.SetEtiketten(c.NewEtiketten)
		}
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	for _, z := range zettels {
		if z, err = u.StoreObjekten().CreateOrUpdate(
			z,
			z.GetKennung(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
