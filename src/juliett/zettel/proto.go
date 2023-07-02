package zettel

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/india/konfig"
)

type ProtoZettel struct {
	Metadatei sku.Metadatei
}

func MakeProtoZettel(k konfig.Compiled) (p ProtoZettel) {
	errors.TodoP1("modify konfig to keep etiketten set")

	p.Metadatei.Typ = k.GetErworben().DefaultTyp

	todo.Decide("should this be set to default etiketten?")
	p.Metadatei.Etiketten = kennung.MakeEtikettSet()

	return
}

func MakeEmptyProtoZettel() ProtoZettel {
	return ProtoZettel{
		Metadatei: sku.Metadatei{
			Etiketten: kennung.MakeEtikettSet(),
		},
	}
}

func (pz *ProtoZettel) AddToFlagSet(f *flag.FlagSet) {
	pz.Metadatei.AddToFlagSet(f)
}

func (pz ProtoZettel) Equals(z sku.Metadatei) (ok bool) {
	var okTyp, okMet bool

	if !kennung.IsEmpty(pz.Metadatei.Typ) &&
		pz.Metadatei.Typ.Equals(z.GetTyp()) {
		okTyp = true
	}

	if pz.Metadatei.Equals(z) {
		okMet = true
	}

	ok = okTyp && okMet

	return
}

func (pz ProtoZettel) Make() (z *sku.Metadatei) {
	todo.Change("add typ")
	todo.Change("add Bezeichnung")
	z = &sku.Metadatei{
		Etiketten: kennung.MakeEtikettSet(),
	}

	pz.Apply(z)

	return
}

func (pz ProtoZettel) Apply(ml metadatei.MetadateiLike) (ok bool) {
	z := ml.GetMetadatei()

	if kennung.IsEmpty(z.GetTyp()) &&
		!kennung.IsEmpty(pz.Metadatei.Typ) &&
		!z.GetTyp().Equals(pz.Metadatei.Typ) {
		ok = true
		z.Typ = pz.Metadatei.Typ
	}

	if pz.Metadatei.Bezeichnung.WasSet() &&
		!z.Bezeichnung.Equals(pz.Metadatei.Bezeichnung) {
		ok = true
		z.Bezeichnung = pz.Metadatei.Bezeichnung
	}

	if pz.Metadatei.Etiketten.Len() > 0 {
		ok = true
	}

	mes := z.Etiketten.MutableClone()
	pz.Metadatei.Etiketten.Each(mes.Add)
	z.Etiketten = mes.ImmutableClone()

	ml.SetMetadatei(z)

	return
}

func (pz ProtoZettel) ApplyWithAkteFD(
	ml metadatei.MetadateiLike,
	akteFD kennung.FD,
) (err error) {
	z := ml.GetMetadatei()

	if kennung.IsEmpty(z.GetTyp()) &&
		!kennung.IsEmpty(pz.Metadatei.Typ) &&
		!z.GetTyp().Equals(pz.Metadatei.Typ) {
		z.Typ = pz.Metadatei.Typ
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
		!z.Bezeichnung.Equals(pz.Metadatei.Bezeichnung) {
		bez = pz.Metadatei.Bezeichnung.String()
	}

	if err = z.Bezeichnung.Set(bez); err != nil {
		err = errors.Wrap(err)
		return
	}

	mes := z.Etiketten.MutableClone()
	pz.Metadatei.Etiketten.Each(mes.Add)
	z.Etiketten = mes.ImmutableClone()

	ml.SetMetadatei(z)

	return
}
