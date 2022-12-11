package commands

import (
	"flag"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/papa/umwelt"
)

type CheckinAkte struct {
	Delete       bool
	NewEtiketten kennung.EtikettSet
}

func init() {
	registerCommand(
		"checkin-akte",
		func(f *flag.FlagSet) Command {
			c := &CheckinAkte{
				NewEtiketten: kennung.MakeEtikettSet(),
			}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.Var(&c.NewEtiketten, "new-etiketten", "comma-separated etiketten (will replace existing Etiketten)")

			return c
		},
	)
}

func (c CheckinAkte) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args)%2 != 0 {
		err = errors.Errorf("arguments must come in pairs of hinweis and akte path")
		return
	}

	type externalAktePair struct {
		hinweis.Hinweis
		path string
	}

	pairs := make([]externalAktePair, len(args)/2)

	// transform args into pairs of hinweis and filepaths
	for i, p := range pairs {
		hs := args[i*2]
		ap := args[(i*2)+1]

		if p.Hinweis, err = hinweis.Make(hs); err != nil {
			err = errors.Wrap(err)
			return
		}

		p.path = ap
		pairs[i] = p
	}

	zettels := make([]zettel.Transacted, len(pairs))
	errors.Log().PrintDebug(pairs)

	// iterate through pairs and read current zettel
	for i, p := range pairs {
		if zettels[i], err = u.StoreObjekten().Zettel().ReadHinweisSchwanzen(p.Hinweis); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for i, p := range pairs {
		var ow sha.WriteCloser

		if ow, err = u.StoreObjekten().AkteWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}

		var f *os.File

		if f, err = files.Open(p.path); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, f.Close)

		if _, err = io.Copy(ow, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = ow.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if zettels[i], err = u.StoreObjekten().Zettel().ReadHinweisSchwanzen(p.Hinweis); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i].Objekte.Akte = ow.Sha()

		if c.NewEtiketten.Len() > 0 {
			zettels[i].Objekte.Etiketten = c.NewEtiketten
		}
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	for _, z := range zettels {
		if z, err = u.StoreObjekten().Zettel().Update(
			&z.Objekte,
			&z.Sku.Kennung,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
