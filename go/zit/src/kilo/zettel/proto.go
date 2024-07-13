package zettel

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
)

type ProtoZettel struct {
	Metadatei object_metadata.Metadatei
}

func MakeProtoZettel(k *konfig.Compiled) (p ProtoZettel) {
	errors.TodoP1("modify konfig to keep etiketten set")

	p.Metadatei.Typ = k.GetErworben().Defaults.Typ
	p.Metadatei.SetEtiketten(k.DefaultEtiketten)

	return
}

func MakeEmptyProtoZettel() ProtoZettel {
	return ProtoZettel{}
}

func (pz *ProtoZettel) AddToFlagSet(f *flag.FlagSet) {
	pz.Metadatei.AddToFlagSet(f)
}

func (pz ProtoZettel) Equals(z *object_metadata.Metadatei) (ok bool) {
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

func (pz ProtoZettel) Make() (z *sku.Transacted) {
	todo.Change("add typ")
	todo.Change("add Bezeichnung")
	z = sku.GetTransactedPool().Get()

	pz.Apply(z, genres.Zettel)

	return
}

func (pz ProtoZettel) Apply(
	ml object_metadata.MetadateiLike,
	gg interfaces.GenreGetter,
) (ok bool) {
	z := ml.GetMetadatei()

	g := gg.GetGenre()
	ui.Log().Print(ml, g)

	switch g {
	case genres.Zettel, genres.Unknown:
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

func (pz ProtoZettel) ApplyWithBlobFD(
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
