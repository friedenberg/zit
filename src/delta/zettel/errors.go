package zettel

import (
	"fmt"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type ErrHasInlineAkteAndFilePath struct {
	FilePath string
	Zettel
	sha.Sha
	AkteWriterFactory
}

func (e ErrHasInlineAkteAndFilePath) Error() string {
	return fmt.Sprintf(
		"zettel text has both inline akte and filepath: %q",
		e.FilePath,
	)
}

func (e ErrHasInlineAkteAndFilePath) Recover() (z Zettel, err error) {
	if e.AkteWriterFactory == nil {
		err = errors.Errorf("akte writer factory is nil")
		return
	}

	var akteWriter sha.WriteCloser

	if akteWriter, err = e.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.Open(e.FilePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	z = e.Zettel
	z.Akte = akteWriter.Sha()

	return
}
