package etiketten_path

import (
	"slices"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/catgut"
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
func (s SliceEtikettWithParents) ContainsKennungEtikett(k *kennung.Kennung2) (int, bool) {
	e := k.PartsStrings().Right
	offset := 0

	if k.IsVirtual() {
		percent := catgut.GetPool().Get()
		defer catgut.GetPool().Put(percent)

		percent.Set("%")

		loc, ok := slices.BinarySearchFunc(
			s,
			percent,
			func(ewp EtikettWithParents, e *Etikett) int {
				cmp := ewp.Etikett.ComparePartial(e)
				return cmp
			},
		)

		if !ok {
			return loc, ok
		}

		offset = percent.Len()
		s = s[loc:]
	}

	return slices.BinarySearchFunc(
		s,
		e,
		func(ewp EtikettWithParents, e *Etikett) int {
			cmp := catgut.CompareUTF8Bytes(ewp.Etikett.Bytes()[offset:], e.Bytes(), true)
			return cmp
		},
	)
}

func (s SliceEtikettWithParents) ContainsEtikett(e *Etikett) (int, bool) {
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
	i, ok := s.ContainsEtikett(e)

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

// TODO return success
func (s *SliceEtikettWithParents) Add(e1 *Etikett, p *Path) (err error) {
	var e *Etikett

	if e, err = e1.Clone(); err != nil {
		err = errors.Wrap(err)
		return
	}

	idx, ok := s.ContainsEtikett(e)

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

// TODO return success
func (s *SliceEtikettWithParents) Remove(e1 *Etikett) (err error) {
	var e *Etikett

	if e, err = e1.Clone(); err != nil {
		err = errors.Wrap(err)
		return
	}

	idx, ok := s.ContainsEtikett(e)

	if !ok {
		return
	}

	*s = slices.Delete(*s, idx, idx+1)

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
