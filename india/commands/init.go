package commands

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
)

type Init struct {
}

func init() {
	registerCommand(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{}

			return c
		},
	)
}

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

	c.mkdirAll(base, "Objekte", "Akte")
	c.mkdirAll(base, "Objekte", "Zettel")

	// c.mkdirAll(base, "Etikett-Zettel")
	c.mkdirAll(base, "Hinweis")
	c.mkdirAll(base, "Zettel-Hinweis")
	c.mkdirAll(base, "Akte-Zettel")
	c.mkdirAll(base, "Verloren+Gefunden")

	c.mkdirAll(base, "Kennung")
	c.writeFile(u.DirZit("Kennung", "Counter"), "0")

	if _, err = _AgeGenerate(base); err != nil {
		err = _Error(err)
		return
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
