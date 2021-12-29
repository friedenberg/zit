package akten

import "fmt"

type DuplicateAkteError struct {
	ShaOldAkte, ShaNewZettel, ShaOldZettel _Sha
}

func (e DuplicateAkteError) Error() string {
	return fmt.Sprintf(
		"already have a zettel for akte:\n      akte: '%s'\nold zettel: '%s'\nnew zettel: '%s'",
		e.ShaOldAkte,
		e.ShaOldZettel,
		e.ShaNewZettel,
	)
}
