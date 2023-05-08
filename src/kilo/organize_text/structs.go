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
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type obj struct {
	Kennung     kennung.Kennung
	Bezeichnung bezeichnung.Bezeichnung
	IsNew       bool
}

func makeObj(
	named metadatei.WithKennung,
) (z obj, err error) {
	errors.TodoP4("add bez in a better way")
	z = obj{
		Kennung:     named.GetKennung(),
		Bezeichnung: bezeichnung.Make(named.GetMetadatei().Description()),
	}

	return
}

func (a obj) Len() int {
	return len(a.Kennung.String())
}

func (a obj) LenKopfUndSchwanz() (int, int) {
	kopf, schwanz := a.KopfUndSchwanz()

	return len(kopf), len(schwanz)
}

func (a obj) KopfUndSchwanz() (kopf, schwanz string) {
	parts := a.Kennung.Parts()
	kopf = parts[0]
	schwanz = parts[2]

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
		if out[i].Bezeichnung.IsEmpty() && out[i].Bezeichnung == out[j].Bezeichnung {
			return out[i].Kennung.String() < out[j].Kennung.String()
		} else {
			return out[i].Bezeichnung.Less(out[j].Bezeichnung)
		}
	})

	return
}
