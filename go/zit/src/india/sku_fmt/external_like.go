package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	ReaderExternalLike = catgut.StringFormatReader[sku.ExternalLike]
	WriterExternalLike = catgut.StringFormatWriter[sku.ExternalLike]

	ExternalLike struct {
		ReaderExternalLike
		WriterExternalLike
	}
)

func MakeExternalLikeCombo(formats map[ids.RepoId]ExternalLike) ExternalLike {
	elc := externalLikeCombo{
		repoIdsToFormats: formats,
	}

	return ExternalLike{
		ReaderExternalLike: elc,
		WriterExternalLike: elc,
	}
}

type externalLikeCombo struct {
	repoIdsToFormats map[ids.RepoId]ExternalLike
}

func (f externalLikeCombo) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	el sku.ExternalLike,
) (n int64, err error) {
	rid := el.GetRepoId()

	w, ok := f.repoIdsToFormats[rid]

	if !ok {
		err = errors.Errorf("no WriterExternalLike for repo id: %s", rid)
		return
	}

	if n, err = w.WriteStringFormat(sw, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f externalLikeCombo) ReadStringFormat(
	rb *catgut.RingBuffer,
	el sku.ExternalLike,
) (n int64, err error) {
	rid := el.GetRepoId()

	r, ok := f.repoIdsToFormats[rid]

	if !ok {
		err = errors.Errorf("no ReaderExternalLike for repo id: %s", rid)
		return
	}

	if n, err = r.ReadStringFormat(rb, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
