package etiketten_path

import (
	"slices"
	"strings"

	"code.linenisgreat.com/zit/src/echo/kennung"
)

type (
	SliceEtikettWithParents []EtikettWithParents
)

func (s *SliceEtikettWithParents) Reset() {
	*s = (*s)[:0]
}

func (s SliceEtikettWithParents) Len() int {
	return len(s)
}

// TODO make less fragile
func (s SliceEtikettWithParents) ContainsEtikett(k *kennung.Kennung2) (int, bool) {
	e := k.PartsStrings().Right

	return slices.BinarySearchFunc(
		s,
		e,
		func(ewp EtikettWithParents, e *Etikett) int {
			cmp := ewp.Etikett.ComparePartial(e)
			return cmp
		},
	)
}

func (s SliceEtikettWithParents) ContainsString(e *Etikett) (int, bool) {
	return slices.BinarySearchFunc(
		s,
		e,
		func(ewp EtikettWithParents, e *Etikett) int {
			cmp := ewp.Etikett.ComparePartial(e)
			return cmp
		},
	)
}

func (s SliceEtikettWithParents) GetMatching(e *Etikett) (matching []EtikettWithParents) {
	i, ok := s.ContainsString(e)

	if !ok {
		return
	}

	for _, ewp := range s[i:] {
		cmp := ewp.ComparePartial(e)

		if cmp != 0 {
			return
		}

		matching = append(matching, ewp)
	}

	return
}

func (s *SliceEtikettWithParents) Add(e *Etikett, p *Path) (err error) {
	idx, ok := s.ContainsString(e)

	var a EtikettWithParents

	if ok {
		a = (*s)[idx]
		a.Parents.AddNonEmptyPath(p)
		(*s)[idx] = a
	} else {
		a = EtikettWithParents{Etikett: e}
		a.Parents.AddNonEmptyPath(p)

		if idx == s.Len() {
			*s = append(*s, a)
		} else {
			*s = slices.Insert(*s, idx, a)
		}
	}

	return
}

func (s SliceEtikettWithParents) StringCommaSeparatedExplicit() string {
	var sb strings.Builder

	first := true

	for _, ewp := range s {
		if ewp.Parents.Len() != 0 {
			continue
		}

		sb.Write(ewp.Etikett.Bytes())

		if !first {
			sb.WriteString(", ")
		}

		first = false
	}

	return sb.String()
}
