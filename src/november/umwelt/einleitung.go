package umwelt

import (
	"bufio"
	"encoding/gob"
	"flag"
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/angeboren"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/typ"
)

type Einleitung struct {
	DisableAge bool
	Yin        string
	Yang       string
	Angeboren  angeboren.Konfig
}

func (e *Einleitung) AddToFlags(f *flag.FlagSet) {
	f.BoolVar(&e.DisableAge, "disable-age", false, "")
	f.StringVar(&e.Yin, "yin", "", "File containing list of Kennung")
	f.StringVar(&e.Yang, "yang", "", "File containing list of Kennung")
	e.Angeboren.AddToFlags(f)
}

func (u *Umwelt) Einleitung(e Einleitung) (err error) {
	s := u.Standort()

	mkdirAll(s.DirKennung())

	if err = readAndTransferLines(e.Yin, s.DirZit("Kennung", "Yin")); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = readAndTransferLines(e.Yang, s.DirZit("Kennung", "Yang")); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, g := range gattung.All() {
		var d string

		if d, err = s.DirObjektenGattung(g); err != nil {
			if gattung.IsErrUnsupportedGattung(err) {
				err = nil
				continue
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		mkdirAll(d)
	}

	mkdirAll(s.DirVerzeichnisse())
	mkdirAll(s.DirVerlorenUndGefunden())

	if !e.DisableAge {
		if _, err = age.Generate(s.FileAge()); err != nil {
			// If the Age file exists, don't do anything and continue init
			if errors.Is(err, os.ErrExist) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}
	}

	writeFile(s.FileKonfigAngeboren(), e.Angeboren)

	if e.Angeboren.UseKonfigErworbenFile {
		writeFile(s.FileKonfigErworben(), "")
	} else {
		writeFile(s.FileKonfigCompiled(), "")
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.TodoP2("determine if this should be an Einleitung option")
	if err = initDefaultTypAndKonfig(u); err != nil {
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

func initDefaultTypAndKonfig(u *Umwelt) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	defaultTyp, defaultTypKennung := typ.Default()

	if _, err = u.StoreObjekten().Typ().ReadOne(defaultTypKennung); err != nil {
		err = nil

		var sh schnittstellen.Sha

		if sh, _, err = u.StoreObjekten().Typ().SaveAkteText(
			*defaultTyp,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defaultTyp.SetAkteSha(sh)

		var defaultTypTransacted *typ.Transacted

		if defaultTypTransacted, err = u.StoreObjekten().Typ().CreateOrUpdate(
			defaultTyp,
			defaultTypKennung,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		u.KonfigPtr().DefaultTyp = *defaultTypTransacted
	}

	{
		defaultKonfig := erworben.Default()

		var sh schnittstellen.Sha

		if sh, _, err = u.StoreObjekten().Konfig().SaveAkteText(
			*defaultKonfig,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defaultKonfig.SetAkteSha(sh)

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

func mkdirAll(elements ...string) {
	err := os.MkdirAll(path.Join(elements...), os.ModeDir|0o755)
	errors.PanicIfError(err)
}

func writeFile(p string, contents any) {
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

func readAndTransferLines(in, out string) (err error) {
	errors.TodoP4("move to user operations")

	if in == "" {
		return
	}

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

		// TODO-P2 sterilize line
		w.WriteString(l)
	}

	return
}
