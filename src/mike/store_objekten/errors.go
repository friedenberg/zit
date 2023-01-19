package store_objekten

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type stringId string

func (s stringId) String() string {
	return string(s)
}

type ErrLockRequired struct {
	Operation string
}

func (e ErrLockRequired) Is(target error) bool {
	_, ok := target.(ErrLockRequired)
	return ok
}

func (e ErrLockRequired) Error() string {
	return fmt.Sprintf(
		"lock required for operation: %q",
		e.Operation,
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

type ErrAkteExists struct {
	Akte sha.Sha
	zettel.MutableSet
}

func (e ErrAkteExists) Is(target error) bool {
	_, ok := target.(ErrAkteExists)
	return ok
}

func (e ErrAkteExists) Error() string {
	return fmt.Sprintf(
		"zettelen already exist with akte:\n%s\n%v",
		e.Akte,
		e.MutableSet,
	)
}
