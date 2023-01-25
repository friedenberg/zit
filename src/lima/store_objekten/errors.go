package store_objekten

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

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
