package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

type zettel struct {
	Hinweis     string
	Bezeichnung string
}

func makeZettel(named zettel_named.Zettel, ha HinweisAbbr) (z zettel, err error) {
	h := named.Hinweis

	if ha != nil {
		if h, err = ha.AbbreviateHinweis(h); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	z = zettel{
		Hinweis:     h.String(),
		Bezeichnung: named.Stored.Zettel.Description(),
	}

	return
}

func (z zettel) String() string {
	return fmt.Sprintf("- [%s] %s", z.Hinweis, z.Bezeichnung)
}

func (z zettel) HinweisAligned(kopf, schwanz int) string {
	parts := strings.Split(z.Hinweis, "/")

	diffKopf := kopf - len(parts[0])
	if diffKopf > 0 {
		parts[0] = strings.Repeat(" ", diffKopf) + parts[0]
	}

	diffSchwanz := schwanz - len(parts[1])
	if diffSchwanz > 0 {
		parts[1] = parts[1] + strings.Repeat(" ", diffSchwanz)
	}

	return fmt.Sprintf("%s/%s", parts[0], parts[1])
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

	z.Hinweis = strings.TrimSpace(remaining[:idx])

	//no bezeichnung
	if idx+2 > len(remaining)-1 {
		return
	}

	remaining = remaining[idx+2:]

	z.Bezeichnung = remaining

	return
}
