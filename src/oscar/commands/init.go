package commands

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/mike/umwelt"
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

func (c Init) Run(u *umwelt.Umwelt, args ...string) (err error) {
	s := u.Standort()
	base := s.DirZit()

	c.mkdirAll(base, "bin")

	c.mkdirAll(s.DirObjektenAkten())
	c.mkdirAll(s.DirObjektenZettelen())
	c.mkdirAll(s.DirObjektenTransaktion())

	c.mkdirAll(s.DirVerlorenUndGefunden())

	c.mkdirAll(s.DirKennung())
	c.writeFile(s.DirZit("Kennung", "Counter"), "0")

	if !c.DisableAge {
		if _, err = age.Generate(s.FileAge()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	c.writeFile(s.DirZit("Konfig"), "")

	if err = c.populateYinIfNecessary(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.populateYangIfNecessary(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer u.Unlock()

	if err = u.StoreObjekten().Reindex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Init) populateYinIfNecessary(s standort.Standort) (err error) {
	if c.Yin == "" {
		return
	}

	err = c.readAndTransferLines(c.Yin, s.DirZit("Kennung", "Yin"))

	return
}

func (c Init) populateYangIfNecessary(s standort.Standort) (err error) {
	if c.Yang == "" {
		return
	}

	err = c.readAndTransferLines(c.Yang, s.DirZit("Kennung", "Yang"))

	return
}

// TODO move to user operations
func (c Init) readAndTransferLines(in, out string) (err error) {
	var fi, fo *os.File

	if fi, err = files.Open(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, fi.Close)

	if fo, err = files.Create(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, fo.Close)

	r := bufio.NewReader(fi)
	w := bufio.NewWriter(fo)

	defer errors.Deferred(&err, w.Flush)

	for {
		var l string
		l, err = r.ReadString('\n')

		if errors.Is(err, io.EOF) {
			err = nil
			break
		}

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		//TODO sterilize line
		w.WriteString(l)
	}

	return
}

func (c Init) mkdirAll(elements ...string) {
	err := os.MkdirAll(path.Join(elements...), os.ModeDir|0755)
	errors.PanicIfError(err)
}

func (c Init) writeFile(path, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0755)
	errors.PanicIfError(err)
}
