package commands

import (
	"flag"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
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

	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

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
		if zettels[i], err = store.StoreObjekten().Read(p.Hinweis); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	//TODO write new akte object for each and update sha
	for i, p := range pairs {
		var ow sha.WriteCloser

		if ow, err = store.StoreObjekten().AkteWriter(); err != nil {
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

		if zettels[i], err = store.StoreObjekten().Read(p.Hinweis); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i].Named.Stored.Zettel.Akte = ow.Sha()

		if c.NewEtiketten.Len() > 0 {
			zettels[i].Named.Stored.Zettel.Etiketten = c.NewEtiketten
		}
	}

	for _, z := range zettels {
		if z, err = store.StoreObjekten().Update(z.Named.Hinweis, z.Named.Stored.Zettel); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.PrintOutf("%s (akte updated)", z.Named)
	}

	return
}
