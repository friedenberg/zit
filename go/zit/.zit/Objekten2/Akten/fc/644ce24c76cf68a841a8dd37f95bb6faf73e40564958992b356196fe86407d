package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

type BlobStore[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] struct {
	dirLayout dir_layout.DirLayout
	Format[A, APtr]
	resetFunc func(APtr)
}

func MakeBlobStore[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
](
	st dir_layout.DirLayout,
	format Format[A, APtr],
	resetFunc func(APtr),
) (s *BlobStore[A, APtr]) {
	s = &BlobStore[A, APtr]{
		dirLayout: st,
		Format:    format,
		resetFunc: resetFunc,
	}

	return
}

func (s *BlobStore[A, APtr]) GetBlob(
	sh interfaces.Sha,
) (a APtr, err error) {
	var ar interfaces.ShaReadCloser

	if ar, err = s.dirLayout.BlobReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	var a1 A
	a = APtr(&a1)
	s.resetFunc(a)

	if _, err = s.ParseBlob(ar, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	actual := ar.GetShaLike()

	if !actual.EqualsSha(sh) {
		err = errors.Errorf("expected sha %s but got %s", sh, actual)
		return
	}

	return
}

func (s *BlobStore[A, APtr]) PutBlob(a APtr) {
	// TODO-P2 implement pool
}

func (h *BlobStore[A, APtr]) SaveBlobText(
	o APtr,
) (sh interfaces.Sha, n int64, err error) {
	var w sha.WriteCloser

	if w, err = h.dirLayout.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if n, err = h.FormatParsedBlob(w, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = w.GetShaLike()

	return
}
