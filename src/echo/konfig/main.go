package konfig

import (
	"encoding/gob"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type Transacted = objekte.Transacted[Objekte, *Objekte, kennung.Konfig, *kennung.Konfig]

type Konfig struct {
	Cli
	Transacted Transacted
}

func Make(s standort.Standort, kc Cli) (c Konfig, err error) {
	c.Transacted.Objekte.Akte = MakeDefaultCompiled()
	c.Cli = kc

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
