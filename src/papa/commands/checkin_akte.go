package commands

import (
	"flag"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CheckinAkte struct {
	Delete       bool
	NewEtiketten etikett.Set
}

func init() {
	registerCommand(
		"checkin-akte",
		func(f *flag.FlagSet) Command {
			c := &CheckinAkte{
				NewEtiketten: etikett.MakeSet(),
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

	zettels := make([]zettel_transacted.Zettel, len(pairs))
	errors.PrintDebug(pairs)

	// iterate through pairs and read current zettel
	for i, p := range pairs {
		if zettels[i], err = u.StoreObjekten().ReadHinweisSchwanzen(p.Hinweis); err != nil {
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

		defer files.Close(f)

		if _, err = io.Copy(ow, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = ow.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if zettels[i], err = u.StoreObjekten().ReadHinweisSchwanzen(p.Hinweis); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i].Named.Stored.Zettel.Akte = ow.Sha()

		if c.NewEtiketten.Len() > 0 {
			zettels[i].Named.Stored.Zettel.Etiketten = c.NewEtiketten
		}
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer u.Unlock()

	for _, z := range zettels {
		if z, err = u.StoreObjekten().Update(&z.Named); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
