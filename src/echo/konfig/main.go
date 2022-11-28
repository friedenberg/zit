package konfig

import (
	"encoding/gob"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type Objekte struct {
	Sha        sha.Sha
	Akte       Compiled
	tomlKonfig tomlKonfig
}

func (o Objekte) AkteSha() sha.Sha {
	return o.Sha
}

type Transacted = objekte.Transacted2[Objekte, *Objekte, kennung.Konfig, *kennung.Konfig]

type Konfig struct {
	Cli
	Transacted Transacted
}

func Make(s standort.Standort, kc Cli) (c Konfig, err error) {
	c.Transacted.Objekte.Akte = MakeDefaultCompiled()
	c.Cli = kc
	// c = DefaultKonfig()

	var f *os.File

	if f, err = files.Open(s.FileKonfigToml()); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}

		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	format := MakeFormatText(metadatei_io.NopAkteFactory())

	if _, err = format.ReadFormat(
		f,
		&c.Transacted.Objekte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.tryReadTransacted(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Konfig) tryReadTransacted(s standort.Standort) (err error) {
	var f *os.File

	if f, err = files.Open(s.FileKonfigCompiled()); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewDecoder(f)

	if err = dec.Decode(&a.Transacted); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (a *Objekte) Equals(b *Objekte) bool {
	panic("TODO not implemented")
	// return false
}

func (a *Objekte) Reset(b *Objekte) {
	panic("TODO not implemented")
	// return false
}

func (c Objekte) Gattung() gattung.Gattung {
	return gattung.Konfig
}
