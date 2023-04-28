package organize_text

import (
	"fmt"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type obj struct {
	kennung.Hinweis
	bezeichnung.Bezeichnung
}

func makeObj(
	named *zettel.Transacted,
	ha schnittstellen.FuncAbbreviateKorper,
) (z obj, err error) {
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
	z = obj{
		Hinweis:     h,
		Bezeichnung: bezeichnung.Make(named.GetMetadatei().Description()),
	}

	return
}

func (a obj) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a obj) Equals(b obj) bool {
	if !a.Hinweis.Equals(b.Hinweis) {
		return false
	}

	if !a.Bezeichnung.Equals(b.Bezeichnung) {
		return false
	}

	return true
}

func (z obj) String() string {
	return fmt.Sprintf("- [%s] %s", z.Hinweis, z.Bezeichnung)
}

func (z *obj) Set(v string) (err error) {
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

	// no bezeichnung
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

func sortObjSet(
	s schnittstellen.MutableSet[obj],
) (out []obj) {
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

type newObj struct {
	bezeichnung.Bezeichnung
}

func (a newObj) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a newObj) Equals(b newObj) bool {
	if !a.Bezeichnung.Equals(b.Bezeichnung) {
		return false
	}

	return true
}

func (z *newObj) Set(v string) (err error) {
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
	s schnittstellen.MutableSet[newObj],
) (sorted []newObj) {
	sorted = s.Elements()

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Bezeichnung.Less(sorted[j].Bezeichnung)
	})

	return
}
