package sku

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
)

func MakeProto(defaults config_mutable_blobs.Defaults) (p Proto) {
	ui.TodoP1("modify konfig to keep etiketten set")

	p.Metadata.Type = defaults.GetType()
	p.Metadata.SetTags(ids.MakeTagSet(defaults.GetTags()...))

	return
}

type Proto struct {
	object_metadata.Metadata
}

func (pz *Proto) SetFlagSet(f *flag.FlagSet) {
	pz.Metadata.SetFlagSet(f)
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

func (pz Proto) Make() (z *Transacted) {
	todo.Change("add type")
	todo.Change("add description")
	z = GetTransactedPool().Get()

	pz.Apply(z, genres.Zettel)

	return
}

func (proto Proto) ApplyType(
	metadataLike object_metadata.MetadataLike,
	genreGetter interfaces.GenreGetter,
) (ok bool) {
	metadata := metadataLike.GetMetadata()

	g := genreGetter.GetGenre()
	ui.Log().Print(metadataLike, g)

	switch g {
	case genres.Zettel, genres.None:
		if ids.IsEmpty(metadata.GetType()) &&
			!ids.IsEmpty(proto.Metadata.Type) &&
			!metadata.GetType().Equals(proto.Metadata.Type) {
			ok = true
			metadata.Type = proto.Metadata.Type
		}
	}

	return
}

func (proto Proto) Apply(
	metadataLike object_metadata.MetadataLike,
	genreGetter interfaces.GenreGetter,
) (ok bool) {
	metadata := metadataLike.GetMetadata()

	if proto.ApplyType(metadataLike, genreGetter) {
		ok = true
	}

	if proto.Metadata.Description.WasSet() &&
		!metadata.Description.Equals(proto.Metadata.Description) {
		ok = true
		metadata.Description = proto.Metadata.Description
	}

	if proto.Metadata.GetTags().Len() > 0 {
		ok = true
	}

	errors.PanicIfError(proto.Metadata.GetTags().EachPtr(metadata.AddTagPtr))

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

	desc := blobFD.FileNameSansExt()

	if pz.Metadata.Description.WasSet() &&
		!z.Description.Equals(pz.Metadata.Description) {
		desc = pz.Metadata.Description.String()
	}

	if err = z.Description.Set(desc); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.PanicIfError(pz.Metadata.GetTags().EachPtr(z.AddTagPtr))

	return
}
