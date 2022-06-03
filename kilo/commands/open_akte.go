package commands

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/exec"
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
	var hins []_Hinweis
	var shas []_Sha

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

	for i, sha := range shas {
		func(sha _Sha) {
			var z _NamedZettel

			if z, err = store.Zettels().Read(sha); err != nil {
				err = errors.Error(err)
				return
			}

			shaAkte := z.Zettel.Akte
			p := store.DirZit("Objekte", "Akte")

			var f *os.File

			var filename string

			if filename, err = _IdMakeDirNecessary(hins[i], dir); err != nil {
				err = errors.Error(err)
				return
			}

			filename = filename + "." + z.Zettel.AkteExt.String()

			if f, err = _Create(filename); err != nil {
				err = errors.Error(err)
				return
			}

			defer _Close(f)

			files[i] = f.Name()

			if err = _ObjekteRead(f, store.Age(), _IdPath(shaAkte, p)); err != nil {
				err = errors.Error(err)
				return
			}
		}(sha)
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
