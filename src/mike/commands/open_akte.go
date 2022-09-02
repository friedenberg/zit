package commands

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/src/alfa/exec"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type OpenAkte struct {
}

func init() {
	registerCommand(
		"open-akte",
		func(f *flag.FlagSet) Command {
			c := &OpenAkte{}

			return commandWithLockedStore{commandWithHinweisen{c}}
		},
	)
}

func (c OpenAkte) RunWithHinweisen(store store_with_lock.Store, hins ...hinweis.Hinweis) (err error) {
	files := make([]string, len(hins))

	dir, err := ioutil.TempDir("", "")

	if err != nil {
		err = errors.Error(err)
		return
	}

	for i, h := range hins {
		func(h hinweis.Hinweis) {
			var tz zettel_transacted.Zettel

			if tz, err = store.StoreObjekten().Read(h); err != nil {
				err = errors.Error(err)
				return
			}

			shaAkte := tz.Named.Stored.Zettel.Akte

			var f *os.File

			var filename string

			if filename, err = id.MakeDirIfNecessary(hins[i], dir); err != nil {
				err = errors.Error(err)
				return
			}

			filename = filename + "." + tz.Named.Stored.Zettel.Typ.String()

			if f, err = open_file_guard.Create(filename); err != nil {
				err = errors.Error(err)
				return
			}

			defer open_file_guard.Close(f)

			files[i] = f.Name()

			var r io.ReadCloser

			if r, err = store.StoreObjekten().AkteReader(shaAkte); err != nil {
				err = errors.Error(err)
				return
			}

			defer r.Close()

			if _, err = io.Copy(f, r); err != nil {
				err = errors.Error(err)
				return
			}
		}(h)
	}

	cmd := exec.ExecCommand(
		"open",
		files,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = errors.Errorf("opening files ('%q'): %s", files, output)
		return
	}

	return
}
