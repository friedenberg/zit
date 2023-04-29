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
	Kennung     kennung.IdLike
	Bezeichnung bezeichnung.Bezeichnung
	IsNew       bool
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
		Kennung:     h,
		Bezeichnung: bezeichnung.Make(named.GetMetadatei().Description()),
	}

	return
}

func (a obj) Aligned(maxKopf, maxSchwanz int) (v string) {
	if h, ok := a.Hinweis(); ok {
		v = kennung.Aligned(h, maxKopf, maxSchwanz)
	} else if a.Kennung != nil {
		errors.TodoP1("implement alignment for non hinweis kennung")
		v = a.Kennung.String()
	} else {
		panic("kennung was nil")
	}

	return
}

func (a obj) LenKopfUndSchwanz() (int, int) {
	kopf, schwanz := a.KopfUndSchwanz()

	return len(kopf), len(schwanz)
}

func (a obj) KopfUndSchwanz() (kopf, schwanz string) {
	if h, ok := a.Hinweis(); ok {
		kopf = h.Kopf()
		schwanz = h.Schwanz()
	} else {
		schwanz = a.Kennung.String()
	}

	return
}

func (a obj) Hinweis() (h kennung.Hinweis, ok bool) {
	h, ok = a.Kennung.(kennung.Hinweis)
	return
}

func (a obj) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a obj) Equals(b obj) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.Bezeichnung.Equals(b.Bezeichnung) {
		return false
	}

	return true
}

func (z obj) AlignedString(maxKopf, maxSchwanz int) string {
	return fmt.Sprintf("- [%s] %s", z.Kennung, z.Bezeichnung)
}

func (z obj) String() string {
	return fmt.Sprintf("- [%s] %s", z.Kennung, z.Bezeichnung)
}

func (z *obj) setExistingObj(v string) (err error) {
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

	if z.Kennung, err = kennung.Make(
		strings.TrimSpace(remaining[:idx]),
	); err != nil {
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

func (z *obj) setNewObj(v string) (err error) {
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

	z.IsNew = true

	return
}

func sortObjSet(
	s schnittstellen.MutableSet[obj],
) (out []obj) {
	out = s.Elements()

	sort.Slice(out, func(i, j int) bool {
		if out[i].Bezeichnung == out[j].Bezeichnung {
			return out[i].Kennung.String() < out[j].Kennung.String()
		} else {
			return out[i].Bezeichnung.Less(out[j].Bezeichnung)
		}
	})

	return
}
