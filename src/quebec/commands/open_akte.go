package commands

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/exec"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type OpenAkte struct{}

func init() {
	registerCommand(
		"open-akte",
		func(f *flag.FlagSet) Command {
			c := &OpenAkte{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c OpenAkte) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
	)
}

func (c OpenAkte) RunWithIds(store *umwelt.Umwelt, is kennung.Set) (err error) {
	hins := is.Hinweisen.ImmutableClone()
	paths := make([]string, hins.Len())

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for i, h := range hins.Elements() {
		func(h kennung.Hinweis) {
			var tz *zettel.Transacted

			if tz, err = store.StoreObjekten().Zettel().ReadOne(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			shaAkte := tz.Objekte.Akte

			var f *os.File

			var filename string

			if filename, err = id.MakeDirIfNecessary(h, dir); err != nil {
				err = errors.Wrap(err)
				return
			}

			filename = filename + "." + tz.Objekte.Typ.String()

			if f, err = files.Create(filename); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, f.Close)

			paths[i] = f.Name()

			var r io.ReadCloser

			if r, err = store.StoreObjekten().AkteReader(shaAkte); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, r.Close)

			if _, err = io.Copy(f, r); err != nil {
				err = errors.Wrap(err)
				return
			}
		}(h)
	}

	cmd := exec.ExecCommand(
		"open",
		paths,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Errorf("opening files ('%q'): %s", paths, output)
		return
	}

	return
}
