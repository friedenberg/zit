package commands

import (
	"flag"
	"io"
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/age_io"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
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
		err = errors.Error(err)
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

		if p.Hinweis, err = hinweis.MakeBlindHinweis(hs); err != nil {
			err = errors.Error(err)
			return
		}

		p.path = ap
		pairs[i] = p
	}

	zettels := make([]stored_zettel.Transacted, len(pairs))
	logz.PrintDebug(pairs)

	// iterate through pairs and read current zettel
	for i, p := range pairs {
		if zettels[i], err = store.Zettels().Read(p.Hinweis); err != nil {
			err = errors.Error(err)
			return
		}
	}

	//TODO write new akte object for each and update sha
	for i, p := range pairs {
		var ow age_io.Writer

		if ow, err = store.Zettels().AkteWriter(); err != nil {
			err = errors.Error(err)
			return
		}

		var f *os.File

		if f, err = open_file_guard.Open(p.path); err != nil {
			err = errors.Error(err)
			return
		}

		defer open_file_guard.Close(f)

		if _, err = io.Copy(ow, f); err != nil {
			err = errors.Error(err)
			return
		}

		if err = ow.Close(); err != nil {
			err = errors.Error(err)
			return
		}

		if zettels[i], err = store.Zettels().Read(p.Hinweis); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i].Zettel.Akte = ow.Sha()

		if c.NewEtiketten.Len() > 0 {
			zettels[i].Zettel.Etiketten = c.NewEtiketten
		}
	}

	for _, z := range zettels {
		if z, err = store.Zettels().Update(z.Hinweis, z.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		stdprinter.Outf("%s (akte updated)", z.Named)
	}

	return
}
