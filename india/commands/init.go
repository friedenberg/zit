package commands

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
)

type Init struct {
	DisableAge bool
}

func init() {
	registerCommand(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{}

			f.BoolVar(&c.DisableAge, "disable-age", false, "")

			return c
		},
	)
}

//TODO use args for kennung files
func (c Init) Run(u _Umwelt, args ...string) (err error) {
	base := u.DirZit()

	defer func() {
		r := recover()

		if r == nil {
			return
		}

		if e, ok := r.(error); ok {
			err = e
		}

		panic(r)
	}()

	c.mkdirAll(base, "bin")

	c.mkdirAll(u.DirAkte())
	c.mkdirAll(u.DirZettel())

	// c.mkdirAll(base, "Etikett-Zettel")
	c.mkdirAll(u.DirHinweis())
	c.mkdirAll(u.DirZettelHinweis())
	c.mkdirAll(u.DirVerlorenUndGefunden())

	c.mkdirAll(u.DirKennung())
	c.writeFile(u.DirZit("Kennung", "Counter"), "0")

	if !c.DisableAge {
		if _, err = _AgeGenerate(u.FileAge()); err != nil {
			err = _Error(err)
			return
		}
	}

	c.writeFile(u.DirZit("Konfig"), "")

	return
}

func (c Init) mkdirAll(elements ...string) {
	err := os.MkdirAll(path.Join(elements...), os.ModeDir|0755)
	_PanicIfError(err)
}

func (c Init) writeFile(path, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0755)
	_PanicIfError(err)
}
