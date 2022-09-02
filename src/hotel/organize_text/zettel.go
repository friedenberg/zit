package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/lima/zettel_named"
)

type zettel struct {
	Hinweis     string
	Bezeichnung string
}

func makeZettel(named zettel_named.Zettel) zettel {
	return zettel{
		Hinweis:     named.Hinweis.String(),
		Bezeichnung: named.Stored.Zettel.Description(),
	}
}

func (z zettel) String() string {
	return fmt.Sprintf("- [%s] %s", z.Hinweis, z.Bezeichnung)
}

func (z *zettel) Set(v string) (err error) {
	remaining := v

	if len(remaining) < 3 {
		err = errors.Errorf("expected at least 3 characters")
		return
	}

	if remaining[:3] != "- [" {
		err = errors.Errorf("expected '- [', but got '%s'", remaining[:3])
		return
	}

	remaining = remaining[3:]

	idx := -1

	if idx = strings.Index(remaining, "]"); idx == -1 {
		err = errors.Errorf("expected ']' after hinweis, but not found")
		return
	}

	z.Hinweis = remaining[:idx]

	//no bezeichnung
	if idx+2 > len(remaining)-1 {
		return
	}

	remaining = remaining[idx+2:]

	z.Bezeichnung = remaining

	return
}
