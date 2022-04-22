package zettel_formats

import (
	"fmt"
	"io"
	"os"

	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/echo/zettel"
)

type ErrHasInlineAkteAndFilePath struct {
	FilePath string
	zettel.Zettel
	sha.Sha
	_AkteWriterFactory
}

func (e ErrHasInlineAkteAndFilePath) Error() string {
	return fmt.Sprintf(
		"zettel text has both inline akte and filepath: %q",
		e.FilePath,
	)
}

func (e ErrHasInlineAkteAndFilePath) Recover() (z zettel.Zettel, err error) {
	if e._AkteWriterFactory == nil {
		err = _Errorf("akte writer factory is nil")
		return
	}

	var akteWriter _ObjekteWriter

	if akteWriter, err = e.AkteWriter(); err != nil {
		err = _Error(err)
		return
	}

	var f *os.File

	if f, err = _Open(e.FilePath); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(f)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = _Error(err)
		return
	}

	z = e.Zettel
	z.Akte = akteWriter.Sha()

	return
}
