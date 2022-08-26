package store_objekten

import (
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/akten"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/google/uuid"
)

type stringId string

func (s stringId) String() string {
	return string(s)
}

type ErrChainIndexOutOfBounds struct {
	hinweis.HinweisWithIndex
	ChainLength int
}

func (e ErrChainIndexOutOfBounds) Is(target error) bool {
	_, ok := target.(ErrChainIndexOutOfBounds)
	return ok
}

func (e ErrChainIndexOutOfBounds) Error() string {
	return fmt.Sprintf(
		"chain for %s is %d long, but requested index %d",
		e.HinweisWithIndex.Hinweis,
		e.ChainLength,
		e.HinweisWithIndex.Index,
	)
}

type ErrNotFound struct {
	Id fmt.Stringer
}

func (e ErrNotFound) Is(target error) bool {
	_, ok := target.(ErrNotFound)
	return ok
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("objekte with id '%s' not found", e.Id)
}

type ErrZettelDidNotChangeSinceUpdate struct {
	NamedZettel stored_zettel.Named
}

func (e ErrZettelDidNotChangeSinceUpdate) Error() string {
	return fmt.Sprintf(
		"zettel has not changed: [%s %s]",
		e.NamedZettel.Hinweis,
		e.NamedZettel.Stored.Sha,
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
