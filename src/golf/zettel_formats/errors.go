package zettel_formats

import (
	"fmt"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type ErrHasInlineAkteAndFilePath struct {
	FilePath string
	zettel.Zettel
	sha.Sha
	zettel.AkteWriterFactory
}

func (e ErrHasInlineAkteAndFilePath) Error() string {
	return fmt.Sprintf(
		"zettel text has both inline akte and filepath: %q",
		e.FilePath,
	)
}

func (e ErrHasInlineAkteAndFilePath) Recover() (z zettel.Zettel, err error) {
	if e.AkteWriterFactory == nil {
		err = errors.Errorf("akte writer factory is nil")
		return
	}

	var akteWriter age_io.Writer

	if akteWriter, err = e.AkteWriter(); err != nil {
		err = errors.Error(err)
		return
	}

	var f *os.File

	if f, err = open_file_guard.Open(e.FilePath); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = errors.Error(err)
		return
	}

	z = e.Zettel
	z.Akte = akteWriter.Sha()

	return
}
