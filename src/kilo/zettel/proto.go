package zettel

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

type ProtoZettel struct {
	Typ         kennung.Typ
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
}

func MakeProtoZettel(k konfig_compiled.Compiled) ProtoZettel {
	errors.TodoP1("modify konfig to keep etiketten set")

	return ProtoZettel{
		Typ:       k.DefaultTyp.Sku.Kennung,
		Etiketten: kennung.MakeEtikettSet(),
	}
}

func MakeEmptyProtoZettel() ProtoZettel {
	return ProtoZettel{
		Etiketten: kennung.MakeEtikettSet(),
	}
}

func (pz *ProtoZettel) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&pz.Typ, "typ", "the Typ to use for created or updated Zettelen")
	f.Var(&pz.Bezeichnung, "bezeichnung", "the Bezeichnung to use for created or updated Zettelen")
	f.Var(&pz.Etiketten, "etiketten", "the Etiketten to use for created or updated Zttelen")
}

func (pz ProtoZettel) Equals(z Objekte) (ok bool) {
	var okTyp, okEt, okBez bool

	if !pz.Typ.IsEmpty() && pz.Typ.Equals(&z.Typ) {
		okTyp = true
	}

	if pz.Etiketten.Len() > 0 && pz.Etiketten.Equals(z.Etiketten) {
		okEt = true
	}

	if !pz.Bezeichnung.WasSet() || pz.Bezeichnung.Equals(z.Bezeichnung) {
		okBez = true
	}

	ok = okTyp && okEt && okBez

	return
}

func (pz ProtoZettel) Make() (z *Objekte) {
	z = &Objekte{
		Etiketten: kennung.MakeEtikettSet(),
	}

	pz.Apply(z)

	return
}

func (pz ProtoZettel) Apply(z *Objekte) (ok bool) {
	if z.Typ.IsEmpty() && !pz.Typ.IsEmpty() && !z.Typ.Equals(&pz.Typ) {
		ok = true
		z.Typ = pz.Typ
	}

	if pz.Bezeichnung.WasSet() && !z.Bezeichnung.Equals(pz.Bezeichnung) {
		ok = true
		z.Bezeichnung = pz.Bezeichnung
	}

	if pz.Etiketten.Len() > 0 {
		ok = true
	}

	mes := z.Etiketten.MutableCopy()
	pz.Etiketten.Each(mes.Add)
	z.Etiketten = mes.Copy()

	return
}
