package commands

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/alfa/exec"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/charlie/id"
	"github.com/friedenberg/zit/charlie/open_file_guard"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type OpenAkte struct {
}

func init() {
	registerCommand(
		"open-akte",
		func(f *flag.FlagSet) Command {
			c := &OpenAkte{}

			return commandWithLockedStore{c}
		},
	)
}

func (c OpenAkte) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	var hins []hinweis.Hinweis
	var shas []sha.Sha

	if shas, hins, err = store.Hinweisen().ReadManyStrings(args...); err != nil {
		err = errors.Error(err)
		return
	}

	files := make([]string, len(shas))

	dir, err := ioutil.TempDir("", "")

	if err != nil {
		err = errors.Error(err)
		return
	}

	for i, s := range shas {
		func(s sha.Sha) {
			var tz stored_zettel.Transacted

			if tz, err = store.Zettels().Read(s); err != nil {
				err = errors.Error(err)
				return
			}

			shaAkte := tz.Zettel.Akte

			var f *os.File

			var filename string

			if filename, err = id.MakeDirIfNecessary(hins[i], dir); err != nil {
				err = errors.Error(err)
				return
			}

			filename = filename + "." + tz.Zettel.AkteExt.String()

			if f, err = open_file_guard.Create(filename); err != nil {
				err = errors.Error(err)
				return
			}

			defer open_file_guard.Close(f)

			files[i] = f.Name()

			var r io.ReadCloser

			if r, err = store.Zettels().AkteReader(shaAkte); err != nil {
				err = errors.Error(err)
				return
			}

			defer r.Close()

			if _, err = io.Copy(f, r); err != nil {
				err = errors.Error(err)
				return
			}
		}(s)
	}

	cmd := exec.ExecCommand(
		"open",
		[]string{"-W"},
		files,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = errors.Errorf("opening files ('%q'): %s", files, output)
		return
	}

	return
}
