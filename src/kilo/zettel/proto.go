package zettel

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

type ProtoZettel struct {
	Metadatei metadatei.Metadatei
}

func MakeProtoZettel(k konfig.Compiled) (p ProtoZettel) {
	errors.TodoP1("modify konfig to keep etiketten set")

	p.Metadatei.Typ = k.GetErworben().Defaults.Typ

	return
}

func MakeEmptyProtoZettel() ProtoZettel {
	return ProtoZettel{}
}

func (pz *ProtoZettel) AddToFlagSet(f *flag.FlagSet) {
	pz.Metadatei.AddToFlagSet(f)
}

func (pz ProtoZettel) Equals(z *metadatei.Metadatei) (ok bool) {
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

func (pz ProtoZettel) Make() (z *metadatei.Metadatei) {
	todo.Change("add typ")
	todo.Change("add Bezeichnung")
	z = metadatei.GetPool().Get()

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

	if pz.Metadatei.GetEtiketten().Len() > 0 {
		ok = true
	}

	errors.PanicIfError(pz.Metadatei.GetEtiketten().EachPtr(z.AddEtikettPtr))

	return
}

func (pz ProtoZettel) ApplyWithAkteFD(
	ml metadatei.MetadateiLike,
	akteFD fd.FD,
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

	errors.PanicIfError(pz.Metadatei.GetEtiketten().EachPtr(z.AddEtikettPtr))

	return
}
