package sku

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type Proto struct {
	Metadatei object_metadata.Metadatei
}

func (pz *Proto) AddToFlagSet(f *flag.FlagSet) {
	pz.Metadatei.AddToFlagSet(f)
}

func (pz Proto) Equals(z *object_metadata.Metadatei) (ok bool) {
	var okTyp, okMet bool

	if !ids.IsEmpty(pz.Metadatei.Typ) &&
		pz.Metadatei.Typ.Equals(z.GetTyp()) {
		okTyp = true
	}

	if pz.Metadatei.Equals(z) {
		okMet = true
	}

	ok = okTyp && okMet

	return
}

func (pz Proto) Make() (z *object_metadata.Metadatei) {
	todo.Change("add typ")
	todo.Change("add Bezeichnung")
	z = object_metadata.GetPool().Get()

	pz.Apply(z, genres.Zettel)

	return
}

func (pz Proto) Apply(
	ml object_metadata.MetadateiLike,
	g interfaces.GenreGetter,
) (ok bool) {
	z := ml.GetMetadatei()

	if g.GetGenre() == genres.Zettel {
		if ids.IsEmpty(z.GetTyp()) &&
			!ids.IsEmpty(pz.Metadatei.Typ) &&
			!z.GetTyp().Equals(pz.Metadatei.Typ) {
			ok = true
			z.Typ = pz.Metadatei.Typ
		}
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

func (pz Proto) ApplyWithAkteFD(
	ml object_metadata.MetadateiLike,
	akteFD *fd.FD,
) (err error) {
	z := ml.GetMetadatei()

	if ids.IsEmpty(z.GetTyp()) &&
		!ids.IsEmpty(pz.Metadatei.Typ) &&
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
