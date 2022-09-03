package store_objekten

import (
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/echo/akten"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
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
	NamedZettel zettel_named.Zettel
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
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	p1 = path.Join(p, uuid.NewString())

	if f, err = files.Create(p1); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	e.Out = f
	format := zettel.Text{}

	if _, err = format.WriteTo(e.FormatContextWrite); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type ErrAkteExists struct {
	Akte sha.Sha
	zettel_transacted.Set
}

func (e ErrAkteExists) Is(target error) bool {
	_, ok := target.(ErrAkteExists)
	return ok
}

func (e ErrAkteExists) Error() string {
	return fmt.Sprintf(
		"zettelen already exist with akte:\n%s\n%v",
		e.Akte,
		e.Set,
	)
}
