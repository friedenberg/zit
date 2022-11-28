package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/typ_toml"
)

type Objekte struct {
	Sha        sha.Sha
	Akte       Compiled
	tomlKonfig tomlKonfig
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.Sha = v
}

func (o Objekte) AkteSha() sha.Sha {
	return o.Sha
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

func (k *Objekte) AddTyp(
	t *typ_toml.Typ,
	tk *kennung.Typ,
) {
	//TODO
	// k.Akte.Typen[tk.String()] = *t

	return
}

func (k *Objekte) Recompile() (err error) {
	//TODO
	if k.Akte, err = makeCompiled(k.tomlKonfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
