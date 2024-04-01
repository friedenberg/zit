package metadatei

import (
	"fmt"

	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/fd"
)

func MakeErrHasInlineAkteAndFilePath(
	akteFD *fd.FD,
	sh *sha.Sha,
) (err *ErrHasInlineAkteAndFilePath) {
	err = &ErrHasInlineAkteAndFilePath{}
	err.AkteFD.ResetWith(akteFD)
	err.InlineSha.SetShaLike(sh)
	return
}

type ErrHasInlineAkteAndFilePath struct {
	AkteFD    fd.FD
	InlineSha sha.Sha
}

func (e *ErrHasInlineAkteAndFilePath) Error() string {
	return fmt.Sprintf(
		"text has inline akte and file: \nexternal path: %s\nexternal sha: %s\ninline sha: %s",
		e.AkteFD.GetPath(),
		e.AkteFD.GetShaLike(),
		&e.InlineSha,
	)
}

func MakeErrHasInlineAkteAndMetadateiSha(
	inline, metadatei *sha.Sha,
) (err *ErrHasInlineAkteAndMetadateiSha) {
	err = &ErrHasInlineAkteAndMetadateiSha{}
	err.MetadateiSha.SetShaLike(metadatei)
	err.InlineSha.SetShaLike(inline)
	return
}

type ErrHasInlineAkteAndMetadateiSha struct {
	InlineSha    sha.Sha
	MetadateiSha sha.Sha
}

func (e *ErrHasInlineAkteAndMetadateiSha) Error() string {
	return fmt.Sprintf(
		"text has inline akte and metadatei sha: \ninline sha: %s\n metadatei sha: %s",
		&e.InlineSha,
		&e.MetadateiSha,
	)
}
