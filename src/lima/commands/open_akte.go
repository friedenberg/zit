package commands

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/exec"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type OpenAkte struct {
}

func init() {
	registerCommand(
		"open-akte",
		func(f *flag.FlagSet) Command {
			c := &OpenAkte{}

			return commandWithHinweisen{c}
		},
	)
}

func (c OpenAkte) RunWithHinweisen(store *umwelt.Umwelt, hins ...hinweis.Hinweis) (err error) {
	paths := make([]string, len(hins))

	dir, err := ioutil.TempDir("", "")

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for i, h := range hins {
		func(h hinweis.Hinweis) {
			var tz zettel_transacted.Zettel

			if tz, err = store.StoreObjekten().Read(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			shaAkte := tz.Named.Stored.Zettel.Akte

			var f *os.File

			var filename string

			if filename, err = id.MakeDirIfNecessary(hins[i], dir); err != nil {
				err = errors.Wrap(err)
				return
			}

			filename = filename + "." + tz.Named.Stored.Zettel.Typ.String()

			if f, err = files.Create(filename); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer files.Close(f)

			paths[i] = f.Name()

			var r io.ReadCloser

			if r, err = store.StoreObjekten().AkteReader(shaAkte); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer r.Close()

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
