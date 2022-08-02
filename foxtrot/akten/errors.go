package akten

import (
	"fmt"

	"github.com/friedenberg/zit/bravo/sha"
)

type DuplicateAkteError struct {
	ShaOldAkte, ShaNewZettel, ShaOldZettel sha.Sha
}

func (e DuplicateAkteError) Error() string {
	return fmt.Sprintf(
		"already have a zettel for akte:\n      akte: '%s'\nold zettel: '%s'\nnew zettel: '%s'",
		e.ShaOldAkte,
		e.ShaOldZettel,
		e.ShaNewZettel,
	)
}
