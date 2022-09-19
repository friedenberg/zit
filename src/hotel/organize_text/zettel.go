package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

type zettel struct {
	hinweis.Hinweis
	bezeichnung.Bezeichnung
}

func makeZettel(named zettel_named.Zettel, ha hinweis.Abbr) (z zettel, err error) {
	h := named.Hinweis

	if ha != nil {
		if h, err = ha.AbbreviateHinweis(h); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	z = zettel{
		Hinweis:     h,
    //TODO do this smart
		Bezeichnung: bezeichnung.Bezeichnung(named.Stored.Zettel.Description()),
	}

	return
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

	if err = z.Hinweis.Set(strings.TrimSpace(remaining[:idx])); err != nil {
		err = errors.Wrap(err)
		return
	}

	//no bezeichnung
	if idx+2 > len(remaining)-1 {
		return
	}

	remaining = remaining[idx+2:]

	if err = z.Bezeichnung.Set(remaining); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
