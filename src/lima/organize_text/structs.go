package organize_text

import (
	"fmt"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/hinweis"
	zettel_pkg "github.com/friedenberg/zit/src/kilo/zettel"
)

//   _____    _   _       _
//  |__  /___| |_| |_ ___| |
//    / // _ \ __| __/ _ \ |
//   / /|  __/ |_| ||  __/ |
//  /____\___|\__|\__\___|_|
//

type zettel struct {
	hinweis.Hinweis
	bezeichnung.Bezeichnung
}

func makeZettel(
	named *zettel_pkg.Transacted,
	ha schnittstellen.FuncAbbreviateKorper,
) (z zettel, err error) {
	h := *named.Kennung()

	if ha != nil {
		var v string

		if v, err = ha(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = h.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	errors.TodoP4("add bez in a better way")
	z = zettel{
		Hinweis:     h,
		Bezeichnung: bezeichnung.Make(named.Objekte.Description()),
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

func sortZettelSet(
	s collections.MutableValueSet[zettel, *zettel],
) (out []zettel) {
	out = s.Elements()

	sort.Slice(out, func(i, j int) bool {
		if out[i].Bezeichnung == out[j].Bezeichnung {
			return out[i].Hinweis.Less(out[j].Hinweis)
		} else {
			return out[i].Bezeichnung.Less(out[j].Bezeichnung)
		}
	})

	return
}

//   _   _                 _____    _   _       _
//  | \ | | _____      __ |__  /___| |_| |_ ___| |
//  |  \| |/ _ \ \ /\ / /   / // _ \ __| __/ _ \ |
//  | |\  |  __/\ V  V /   / /|  __/ |_| ||  __/ |
//  |_| \_|\___| \_/\_/   /____\___|\__|\__\___|_|
//

type newZettel struct {
	bezeichnung.Bezeichnung
}

func (z *newZettel) Set(v string) (err error) {
	remaining := v

	if remaining[:2] != "- " {
		err = errors.Errorf("expected '- ', but got '%s'", remaining[:2])
		return
	}

	remaining = remaining[2:]

	if err = z.Bezeichnung.Set(remaining); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func sortNewZettelSet(
	s collections.MutableValueSet[newZettel, *newZettel],
) (sorted []newZettel) {
	sorted = s.Elements()

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Bezeichnung.Less(sorted[j].Bezeichnung)
	})

	return
}
