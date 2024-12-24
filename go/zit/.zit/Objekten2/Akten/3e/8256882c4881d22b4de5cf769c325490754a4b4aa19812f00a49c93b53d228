package object_metadata

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

func MakeErrHasInlineBlobAndFilePath(
	blobFD *fd.FD,
	sh *sha.Sha,
) (err *ErrHasInlineBlobAndFilePath) {
	err = &ErrHasInlineBlobAndFilePath{}
	err.BlobFD.ResetWith(blobFD)
	err.InlineSha.SetShaLike(sh)
	return
}

type ErrHasInlineBlobAndFilePath struct {
	BlobFD    fd.FD
	InlineSha sha.Sha
}

func (e *ErrHasInlineBlobAndFilePath) Error() string {
	return fmt.Sprintf(
		"text has inline blob and file: \nexternal path: %s\nexternal sha: %s\ninline sha: %s",
		e.BlobFD.GetPath(),
		e.BlobFD.GetShaLike(),
		&e.InlineSha,
	)
}

func MakeErrHasInlineBlobAndMetadataSha(
	inline, object_metadata *sha.Sha,
) (err *ErrHasInlineBlobAndMetadataSha) {
	err = &ErrHasInlineBlobAndMetadataSha{}
	err.MetadataSha.SetShaLike(object_metadata)
	err.InlineSha.SetShaLike(inline)
	return
}

type ErrHasInlineBlobAndMetadataSha struct {
	InlineSha   sha.Sha
	MetadataSha sha.Sha
}

func (e *ErrHasInlineBlobAndMetadataSha) Error() string {
	return fmt.Sprintf(
		"text has inline blob and metadata sha: \ninline sha: %s\n metadata sha: %s",
		&e.InlineSha,
		&e.MetadataSha,
	)
}
