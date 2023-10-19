package metadatei

import (
	"fmt"

	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/fd"
)

type ErrHasInlineAkteAndFilePath struct {
	AkteFD    fd.FD
	InlineSha sha.Sha
}

func (e ErrHasInlineAkteAndFilePath) Error() string {
	return fmt.Sprintf(
		"text has inline akte and file: \nexternal path: %s\nexternal sha: %s\ninline sha: %s",
		e.AkteFD.GetPath(),
		e.AkteFD.GetShaLike(),
		e.InlineSha,
	)
}

type ErrHasInlineAkteAndMetadateiSha struct {
	InlineSha    sha.Sha
	MetadateiSha sha.Sha
}

func (e ErrHasInlineAkteAndMetadateiSha) Error() string {
	return fmt.Sprintf(
		"text has inline akte and metadatei sha: \ninline sha: %s\n metadatei sha: %s",
		e.InlineSha,
		e.MetadateiSha,
	)
}
