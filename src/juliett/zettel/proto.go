package zettel

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/india/konfig"
)

type ProtoZettel struct {
	Typ       kennung.Typ
	Metadatei metadatei.Metadatei
}

func MakeProtoZettel(k konfig.Compiled) (p ProtoZettel) {
	errors.TodoP1("modify konfig to keep etiketten set")

	p.Typ = k.GetErworben().DefaultTyp

	todo.Decide("should this be set to default etiketten?")
	p.Metadatei.Etiketten = kennung.MakeEtikettSet()

	return
}

func MakeEmptyProtoZettel() ProtoZettel {
	return ProtoZettel{
		Metadatei: metadatei.Metadatei{
			Etiketten: kennung.MakeEtikettSet(),
		},
	}
}

func (pz *ProtoZettel) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&pz.Typ, "typ", "the Typ to use for created or updated Zettelen")
	pz.Metadatei.AddToFlagSet(f)
}

func (pz ProtoZettel) Equals(z Objekte) (ok bool) {
	var okTyp, okMet bool

	if !pz.Typ.IsEmpty() && pz.Typ.Equals(z.Typ) {
		okTyp = true
	}

	if pz.Metadatei.Equals(z.Metadatei) {
		okMet = true
	}

	ok = okTyp && okMet

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

	if pz.Metadatei.Bezeichnung.WasSet() &&
		!z.Metadatei.Bezeichnung.Equals(pz.Metadatei.Bezeichnung) {
		ok = true
		z.Metadatei.Bezeichnung = pz.Metadatei.Bezeichnung
	}

	if pz.Metadatei.Etiketten.Len() > 0 {
		ok = true
	}

	mes := z.Metadatei.Etiketten.MutableClone()
	pz.Metadatei.Etiketten.Each(mes.Add)
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

	if pz.Metadatei.Bezeichnung.WasSet() &&
		!z.Metadatei.Bezeichnung.Equals(pz.Metadatei.Bezeichnung) {
		bez = pz.Metadatei.Bezeichnung.String()
	}

	if err = z.Metadatei.Bezeichnung.Set(bez); err != nil {
		err = errors.Wrap(err)
		return
	}

	mes := z.Metadatei.Etiketten.MutableClone()
	pz.Metadatei.Etiketten.Each(mes.Add)
	z.Metadatei.Etiketten = mes.ImmutableClone()

	return
}
