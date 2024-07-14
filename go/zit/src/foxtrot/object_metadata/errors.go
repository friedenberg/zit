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

func MakeErrHasInlineBlobAndMetadateiSha(
	inline, object_metadata *sha.Sha,
) (err *ErrHasInlineBlobAndMetadateiSha) {
	err = &ErrHasInlineBlobAndMetadateiSha{}
	err.MetadataSha.SetShaLike(object_metadata)
	err.InlineSha.SetShaLike(inline)
	return
}

type ErrHasInlineBlobAndMetadateiSha struct {
	InlineSha   sha.Sha
	MetadataSha sha.Sha
}

func (e *ErrHasInlineBlobAndMetadateiSha) Error() string {
	return fmt.Sprintf(
		"text has inline blob and metadatei sha: \ninline sha: %s\n metadatei sha: %s",
		&e.InlineSha,
		&e.MetadataSha,
	)
}
