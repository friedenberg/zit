package commands

import (
	"bufio"
	"encoding/gob"
	"flag"
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/angeboren"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Init struct {
	DisableAge bool
	Yin        string
	Yang       string
	Angeboren  angeboren.Konfig
}

func init() {
	registerCommand(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{
				Angeboren: angeboren.Default(),
			}

			f.BoolVar(&c.DisableAge, "disable-age", false, "")
			f.StringVar(&c.Yin, "yin", "", "File containing list of Kennung")
			f.StringVar(&c.Yang, "yang", "", "File containing list of Kennung")
			c.Angeboren.AddToFlags(f)

			return c
		},
	)
}

func (c Init) Run(u *umwelt.Umwelt, args ...string) (err error) {
	s := u.Standort()
	base := s.DirZit()

	c.mkdirAll(base, "bin")

	for _, g := range gattung.All() {
		var d string

		if d, err = s.DirObjektenGattung(g); err != nil {
			if errors.Is(err, gattung.ErrUnsupportedGattung) {
				err = nil
				continue
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		c.mkdirAll(d)
	}

	c.mkdirAll(s.DirVerzeichnisse())
	c.mkdirAll(s.DirVerlorenUndGefunden())

	c.mkdirAll(s.DirKennung())
	c.writeFile(s.DirZit("Kennung", "Counter"), "0")

	if !c.DisableAge {
		if _, err = age.Generate(s.FileAge()); err != nil {
			//If the Age file exists, don't do anything and continue init
			if errors.Is(err, os.ErrExist) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}
	}

	c.writeFile(s.FileKonfigAngeboren(), c.Angeboren)

	if c.Angeboren.UseKonfigErworbenFile {
		c.writeFile(s.FileKonfigErworben(), "")
	} else {
		c.writeFile(s.FileKonfigCompiled(), "")
	}

	//TODO-P2 how to handle re-init for yin and yang?
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

	if err = c.initDefaultTypAndKonfig(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		if err = u.Lock(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, u.Unlock)

		if err = u.StoreObjekten().Reindex(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Init) initDefaultTypAndKonfig(u *umwelt.Umwelt) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	defaultTyp, defaultTypKennung := typ.Default()

	if _, err = u.StoreObjekten().Typ().ReadOne(defaultTypKennung); err != nil {
		err = nil

		if _, err = u.StoreObjekten().Typ().SaveAkteText(
			defaultTyp,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var defaultTypTransacted *typ.Transacted

		if defaultTypTransacted, err = u.StoreObjekten().Typ().CreateOrUpdate(
			defaultTyp,
			defaultTypKennung,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		u.KonfigPtr().AddTyp(defaultTypTransacted)
		u.KonfigPtr().DefaultTyp = *defaultTypTransacted
	}

	{
		defaultKonfig := erworben.Default()

		if _, err = u.StoreObjekten().Konfig().SaveAkteText(
			defaultKonfig,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var defaultKonfigTransacted *erworben.Transacted

		if defaultKonfigTransacted, err = u.StoreObjekten().Konfig().Update(
			defaultKonfig,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		u.KonfigPtr().SetTransacted(defaultKonfigTransacted)
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

// TODO-P4 move to user operations
func (c Init) readAndTransferLines(in, out string) (err error) {
	var fi, fo *os.File

	if fi, err = files.Open(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, fi.Close)

	if fo, err = files.CreateExclusiveWriteOnly(out); err != nil {
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

		//TODO-P2 sterilize line
		w.WriteString(l)
	}

	return
}

func (c Init) mkdirAll(elements ...string) {
	err := os.MkdirAll(path.Join(elements...), os.ModeDir|0755)
	errors.PanicIfError(err)
}

func (c Init) writeFile(p string, contents any) {
	var f *os.File
	var err error

	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
		if errors.IsExist(err) {
			errors.Err().Printf("%s already exists, not overwriting", p)
			err = nil
		} else {
		}

		return
	}

	defer errors.PanicIfError(err)
	defer errors.DeferredCloser(&err, f)

	if s, ok := contents.(string); ok {
		_, err = io.WriteString(f, s)
	} else {
		enc := gob.NewEncoder(f)
		err = enc.Encode(contents)
	}
}
