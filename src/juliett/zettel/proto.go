package zettel

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/india/konfig"
)

type ProtoZettel struct {
	Typ         kennung.Typ
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
}

func MakeProtoZettel(k konfig.Compiled) (p ProtoZettel) {
	errors.TodoP1("modify konfig to keep etiketten set")

	p.Typ = k.GetErworben().DefaultTyp

	todo.Decide("should this be set to default etiketten?")
	p.Etiketten = kennung.MakeEtikettSet()

	return
}

func MakeEmptyProtoZettel() ProtoZettel {
	return ProtoZettel{
		Etiketten: kennung.MakeEtikettSet(),
	}
}

func (pz *ProtoZettel) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&pz.Typ, "typ", "the Typ to use for created or updated Zettelen")
	f.Var(&pz.Bezeichnung, "bezeichnung", "the Bezeichnung to use for created or updated Zettelen")
	f.Var(
		collections.MakeFlagCommasFromExisting(
			collections.SetterPolicyAppend,
			&pz.Etiketten,
		),
		"etiketten",
		"the Etiketten to use for created or updated Zttelen",
	)
}

func (pz ProtoZettel) Equals(z Objekte) (ok bool) {
	var okTyp, okEt, okBez bool

	if !pz.Typ.IsEmpty() && pz.Typ.Equals(z.Typ) {
		okTyp = true
	}

	if pz.Etiketten.Len() > 0 && pz.Etiketten.Equals(z.Metadatei.Etiketten) {
		okEt = true
	}

	if !pz.Bezeichnung.WasSet() || pz.Bezeichnung.Equals(z.Metadatei.Bezeichnung) {
		okBez = true
	}

	ok = okTyp && okEt && okBez

	return
}

func (pz ProtoZettel) Make() (z *Objekte) {
	todo.Change("add typ")
	todo.Change("add Bezeichnung")
	z = &Objekte{
		Metadatei: metadatei.Metadatei{
			Etiketten: kennung.MakeEtikettSet(),
		},
	}

	pz.Apply(z)

	return
}

func (pz ProtoZettel) Apply(z *Objekte) (ok bool) {
	if z.Typ.IsEmpty() && !pz.Typ.IsEmpty() && !z.Typ.Equals(pz.Typ) {
		ok = true
		z.Typ = pz.Typ
	}

	if pz.Bezeichnung.WasSet() && !z.Metadatei.Bezeichnung.Equals(pz.Bezeichnung) {
		ok = true
		z.Metadatei.Bezeichnung = pz.Bezeichnung
	}

	if pz.Etiketten.Len() > 0 {
		ok = true
	}

	mes := z.Metadatei.Etiketten.MutableClone()
	pz.Etiketten.Each(mes.Add)
	z.Metadatei.Etiketten = mes.ImmutableClone()

	return
}

func (pz ProtoZettel) ApplyWithAkteFD(z *Objekte, akteFD kennung.FD) (err error) {
	if z.Typ.IsEmpty() && !pz.Typ.IsEmpty() && !z.Typ.Equals(pz.Typ) {
		z.Typ = pz.Typ
	} else {
		// TODO-P4 use konfig
		ext := akteFD.Ext()

		if ext != "" {
			if err = z.Typ.Set(akteFD.Ext()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	bez := akteFD.FileNameSansExt()

	if pz.Bezeichnung.WasSet() && !z.Metadatei.Bezeichnung.Equals(pz.Bezeichnung) {
		bez = pz.Bezeichnung.String()
	}

	if err = z.Metadatei.Bezeichnung.Set(bez); err != nil {
		err = errors.Wrap(err)
		return
	}

	mes := z.Metadatei.Etiketten.MutableClone()
	pz.Etiketten.Each(mes.Add)
	z.Metadatei.Etiketten = mes.ImmutableClone()

	return
}
