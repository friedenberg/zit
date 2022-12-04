package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/typ_toml"
)

type Objekte struct {
	Sha sha.Sha

	wasCompiled bool
	Akte        Compiled
	Toml        objekteToml
}

type objekteToml struct {
	Sha    sha.Sha
	Konfig tomlKonfig
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.Sha = v
}

func (o Objekte) AkteSha() sha.Sha {
	return o.Sha
}

func (a Objekte) Equals(b *Objekte) bool {
	if b == nil {
		return false
	}

	if !b.wasCompiled {
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		return false
	}

	return true
}

func (a *Objekte) Reset(b *Objekte) {
	if b == nil {
		a.Sha = b.Sha
		a.Akte = b.Akte
		a.Toml = b.Toml
	} else {
		a.Sha = sha.Sha{}
		a.Akte = MakeDefaultCompiled()
		a.Toml = objekteToml{}
	}
}

func (c Objekte) Gattung() gattung.Gattung {
	return gattung.Konfig
}

func (k *Objekte) AddTyp(
	t *typ_toml.Typ,
	tk *kennung.Typ,
) {
	ct := makeCompiledTyp(tk.String())
	ct.Typ.Akte.Apply(t)
	m := k.Akte.Typen.Elements()
	m = append(m, ct)
	k.Akte.Typen = makeCompiledTypSetFromSlice(m)

	return
}

func (k *Objekte) Recompile() (err error) {
	if k.Sha, err = k.Akte.recompile(k.Toml.Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.wasCompiled = true

	return
}
