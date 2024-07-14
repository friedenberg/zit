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
	object_metadata.Metadata
}

func (pz *Proto) AddToFlagSet(f *flag.FlagSet) {
	pz.Metadata.AddToFlagSet(f)
}

func (pz Proto) Equals(z *object_metadata.Metadata) (ok bool) {
	var okTyp, okMet bool

	if !ids.IsEmpty(pz.Metadata.Type) &&
		pz.Metadata.Type.Equals(z.GetType()) {
		okTyp = true
	}

	if pz.Metadata.Equals(z) {
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
			!ids.IsEmpty(pz.Metadata.Type) &&
			!z.GetType().Equals(pz.Metadata.Type) {
			ok = true
			z.Type = pz.Metadata.Type
		}
	}

	if pz.Metadata.Description.WasSet() &&
		!z.Description.Equals(pz.Metadata.Description) {
		ok = true
		z.Description = pz.Metadata.Description
	}

	if pz.Metadata.GetTags().Len() > 0 {
		ok = true
	}

	errors.PanicIfError(pz.Metadata.GetTags().EachPtr(z.AddTagPtr))

	return
}

func (pz Proto) ApplyWithBlobFD(
	ml object_metadata.MetadataLike,
	blobFD *fd.FD,
) (err error) {
	z := ml.GetMetadata()

	if ids.IsEmpty(z.GetType()) &&
		!ids.IsEmpty(pz.Metadata.Type) &&
		!z.GetType().Equals(pz.Metadata.Type) {
		z.Type = pz.Metadata.Type
	} else {
		// TODO-P4 use konfig
		ext := blobFD.Ext()

		if ext != "" {
			if err = z.Type.Set(blobFD.Ext()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	bez := blobFD.FileNameSansExt()

	if pz.Metadata.Description.WasSet() &&
		!z.Description.Equals(pz.Metadata.Description) {
		bez = pz.Metadata.Description.String()
	}

	if err = z.Description.Set(bez); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.PanicIfError(pz.Metadata.GetTags().EachPtr(z.AddTagPtr))

	return
}
