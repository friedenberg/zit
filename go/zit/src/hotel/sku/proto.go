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
	Metadatei object_metadata.Metadata
}

func (pz *Proto) AddToFlagSet(f *flag.FlagSet) {
	pz.Metadatei.AddToFlagSet(f)
}

func (pz Proto) Equals(z *object_metadata.Metadata) (ok bool) {
	var okTyp, okMet bool

	if !ids.IsEmpty(pz.Metadatei.Type) &&
		pz.Metadatei.Type.Equals(z.GetType()) {
		okTyp = true
	}

	if pz.Metadatei.Equals(z) {
		okMet = true
	}

	ok = okTyp && okMet

	return
}

func (pz Proto) Make() (z *object_metadata.Metadata) {
	todo.Change("add typ")
	todo.Change("add Bezeichnung")
	z = object_metadata.GetPool().Get()

	pz.Apply(z, genres.Zettel)

	return
}

func (pz Proto) Apply(
	ml object_metadata.MetadataLike,
	g interfaces.GenreGetter,
) (ok bool) {
	z := ml.GetMetadata()

	if g.GetGenre() == genres.Zettel {
		if ids.IsEmpty(z.GetType()) &&
			!ids.IsEmpty(pz.Metadatei.Type) &&
			!z.GetType().Equals(pz.Metadatei.Type) {
			ok = true
			z.Type = pz.Metadatei.Type
		}
	}

	if pz.Metadatei.Description.WasSet() &&
		!z.Description.Equals(pz.Metadatei.Description) {
		ok = true
		z.Description = pz.Metadatei.Description
	}

	if pz.Metadatei.GetTags().Len() > 0 {
		ok = true
	}

	errors.PanicIfError(pz.Metadatei.GetTags().EachPtr(z.AddTagPtr))

	return
}

func (pz Proto) ApplyWithAkteFD(
	ml object_metadata.MetadataLike,
	akteFD *fd.FD,
) (err error) {
	z := ml.GetMetadata()

	if ids.IsEmpty(z.GetType()) &&
		!ids.IsEmpty(pz.Metadatei.Type) &&
		!z.GetType().Equals(pz.Metadatei.Type) {
		z.Type = pz.Metadatei.Type
	} else {
		// TODO-P4 use konfig
		ext := akteFD.Ext()

		if ext != "" {
			if err = z.Type.Set(akteFD.Ext()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	bez := akteFD.FileNameSansExt()

	if pz.Metadatei.Description.WasSet() &&
		!z.Description.Equals(pz.Metadatei.Description) {
		bez = pz.Metadatei.Description.String()
	}

	if err = z.Description.Set(bez); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.PanicIfError(pz.Metadatei.GetTags().EachPtr(z.AddTagPtr))

	return
}
