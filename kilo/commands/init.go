package commands

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
)

type Init struct {
	DisableAge bool
	Yin        string
	Yang       string
}

func init() {
	registerCommand(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{}

			f.BoolVar(&c.DisableAge, "disable-age", false, "")
			f.StringVar(&c.Yin, "yin", "", "File containing list of Kennung")
			f.StringVar(&c.Yang, "yang", "", "File containing list of Kennung")

			return c
		},
	)
}

func (c Init) Run(u _Umwelt, args ...string) (err error) {
	base := u.DirZit()

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
			err = errors.Error(err)
			return
		}
	}

	c.writeFile(u.DirZit("Konfig"), "")

	if err = c.populateYinIfNecessary(u); err != nil {
		err = errors.Error(err)
		return
	}

	if err = c.populateYangIfNecessary(u); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (c Init) populateYinIfNecessary(u _Umwelt) (err error) {
	if c.Yin == "" {
		return
	}

	err = c.readAndTransferLines(c.Yin, u.DirZit("Kennung", "Yin"))

	return
}

func (c Init) populateYangIfNecessary(u _Umwelt) (err error) {
	if c.Yang == "" {
		return
	}

	err = c.readAndTransferLines(c.Yang, u.DirZit("Kennung", "Yang"))

	return
}

//TODO move to user operations
func (c Init) readAndTransferLines(in, out string) (err error) {
	var fi, fo *os.File

	if fi, err = _Open(in); err != nil {
		err = errors.Error(err)
		return
	}

	defer stdprinter.PanicIfError(fi.Close)

	if fo, err = _Create(out); err != nil {
		err = errors.Error(err)
		return
	}

	defer stdprinter.PanicIfError(fo.Close)

	r := bufio.NewReader(fi)
	w := bufio.NewWriter(fo)

	defer stdprinter.PanicIfError(w.Flush)

	for {
		var l string
		l, err = r.ReadString('\n')

		if errors.Is(err, io.EOF) {
			err = nil
			break
		}

		if err != nil {
			err = errors.Error(err)
			return
		}

		//TODO sterilize line
		w.WriteString(l)
	}

	return
}

func (c Init) mkdirAll(elements ...string) {
	err := os.MkdirAll(path.Join(elements...), os.ModeDir|0755)
	stdprinter.PanicIfError(err)
}

func (c Init) writeFile(path, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0755)
	stdprinter.PanicIfError(err)
}
