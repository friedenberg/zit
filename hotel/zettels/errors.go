package zettels

import (
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

type ErrZettelDidNotChangeSinceUpdate struct {
	NamedZettel _NamedZettel
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

type duplicateAkteError struct {
	_ErrorDuplicateAtke
	_ZettelFormatContextWrite
}

func (e duplicateAkteError) AddToLostAndFound(p string) (p1 string, err error) {
	newEtikett := "zz-akte-" + e.ShaOldAkte.String()

	if err = e.Zettel.Etiketten.AddString(newEtikett); err != nil {
		err = _Error(err)
		return
	}

	var f *os.File

	p1 = path.Join(p, uuid.NewString())

	if f, err = _Create(p1); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(f)

	e.Out = f
	format := _ZettelFormatText{}

	if _, err = format.WriteTo(e._ZettelFormatContextWrite); err != nil {
		err = _Error(err)
		return
	}

	return
}
