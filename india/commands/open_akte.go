package commands

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/alfa/exec"
)

type OpenAkte struct {
}

func init() {
	registerCommand(
		"open-akte",
		func(f *flag.FlagSet) Command {
			c := &OpenAkte{}

			return commandWithZettels{c}
		},
	)
}

func (c OpenAkte) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	var hins []_Hinweis
	var shas []_Sha

	if shas, hins, err = zs.Hinweisen().ReadManyStrings(args...); err != nil {
		err = _Error(err)
		return
	}

	files := make([]string, len(shas))

	dir, err := ioutil.TempDir("", "")

	if err != nil {
		err = _Error(err)
		return
	}

	for i, sha := range shas {
		func(sha _Sha) {
			var z _NamedZettel

			if z, err = zs.Read(sha); err != nil {
				err = _Error(err)
				return
			}

			shaAkte := z.Zettel.Akte
			p := u.DirZit("Objekte", "Akte")

			var f *os.File

			var filename string

			if filename, err = _IdMakeDirNecessary(hins[i], dir); err != nil {
				err = _Error(err)
				return
			}

			filename = filename + "." + z.Zettel.AkteExt.String()

			if f, err = _Create(filename); err != nil {
				err = _Error(err)
				return
			}

			defer _Close(f)

			files[i] = f.Name()

			if err = _ObjekteRead(f, zs.Age(), _IdPath(shaAkte, p)); err != nil {
				err = _Error(err)
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
		err = _Errorf("opening files ('%q'): %s", files, output)
		return
	}

	return
}
