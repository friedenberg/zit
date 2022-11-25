package konfig

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
	"github.com/friedenberg/zit/src/fd"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type Stored = objekte.Stored[Objekte, *Objekte]
type Named = objekte.Named[Objekte, *Objekte, kennung.Konfig, *kennung.Konfig]
type Transacted = objekte.Transacted[Objekte, *Objekte, kennung.Konfig, *kennung.Konfig]

type External struct {
	Named Named
	FD    fd.FD
}

type Objekte struct {
	Sha  sha.Sha
	Akte Konfig
}

type Konfig struct {
	Cli
	tomlKonfig
	Compiled
}

func Make(p string, kc Cli) (c Objekte, err error) {
	c.Akte.Compiled = MakeDefaultCompiled()
	c.Akte.Cli = kc
	// c = DefaultKonfig()

	var f *os.File

	if f, err = files.Open(p); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}

		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	format := MakeFormatText(metadatei_io.NopAkteFactory())

	if _, err = format.ReadFormat(f, &c); err != nil {
		err = errors.Wrap(err)
		return
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
