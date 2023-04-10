package zettel

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha"
)

type externalFile struct {
	Sha  sha.Sha
	Path string
}

type ErrHasInlineAkteAndFilePath struct {
	External  externalFile
	InlineSha sha.Sha
	Objekte
}

func (e ErrHasInlineAkteAndFilePath) Error() string {
	return fmt.Sprintf(
		"zettel text has both inline akte and filepath: \nexternal path: %s\nexternal sha: %s\ninline sha: %s",
		e.External.Path,
		e.External.Sha,
		e.InlineSha,
	)
}

type ErrHasInvalidAkteShaOrFilePath struct {
	Value string
}

func (e ErrHasInvalidAkteShaOrFilePath) Error() string {
	return fmt.Sprintf(
		"zettel text has invalid akte sha or file path: %q",
		e.Value,
	)
}

func (e ErrHasInvalidAkteShaOrFilePath) Is(target error) (ok bool) {
	_, ok = target.(ErrHasInvalidAkteShaOrFilePath)
	return
}
