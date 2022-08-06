package zettels

import (
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/akten"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/google/uuid"
)

type ErrZettelDidNotChangeSinceUpdate struct {
	NamedZettel stored_zettel.Named
}

func (e ErrZettelDidNotChangeSinceUpdate) Error() string {
	return fmt.Sprintf(
		"zettel has not changed: [%s %s]",
		e.NamedZettel.Hinweis,
		e.NamedZettel.Sha,
	)
}

type VerlorenAndGefundenError interface {
	error
	AddToLostAndFound(string) (string, error)
}

type ErrZettelSplitHistory struct {
	hinweis.Hinweis
	ShaA, ShaB sha.Sha
}

func (e ErrZettelSplitHistory) Error() string {
	return fmt.Sprintf(
		"two separate zettels with hinweis:\n%s:\n%s\n%s",
		e.Hinweis,
		e.ShaA,
		e.ShaB,
	)
}

type duplicateAkteError struct {
	akten.DuplicateAkteError
	zettel.FormatContextWrite
}

func (e duplicateAkteError) AddToLostAndFound(p string) (p1 string, err error) {
	newEtikett := "zz-akte-" + e.ShaOldAkte.String()

	if err = e.Zettel.Etiketten.AddString(newEtikett); err != nil {
		err = errors.Error(err)
		return
	}

	var f *os.File

	p1 = path.Join(p, uuid.NewString())

	if f, err = open_file_guard.Create(p1); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	e.Out = f
	format := zettel_formats.Text{}

	if _, err = format.WriteTo(e.FormatContextWrite); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
