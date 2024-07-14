package tag_paths

import (
	"slices"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	TagsWithParentsAndTypes []EtikettWithParentsAndTypes
)

func (s *TagsWithParentsAndTypes) Reset() {
	*s = (*s)[:0]
}

func (s TagsWithParentsAndTypes) Len() int {
	return len(s)
}

func (s TagsWithParentsAndTypes) ContainsKennungEtikett(
	k *ids.ObjectId,
) (int, bool) {
	return s.containsKennungEtikett(k, true)
}

func (s TagsWithParentsAndTypes) ContainsKennungEtikettExact(
	k *ids.ObjectId,
) (int, bool) {
	return s.containsKennungEtikett(k, false)
}

// TODO make less fragile
func (s TagsWithParentsAndTypes) containsKennungEtikett(
	k *ids.ObjectId,
	partial bool,
) (int, bool) {
	e := k.PartsStrings().Right
	offset := 0

	if k.IsVirtual() {
		percent := catgut.GetPool().Get()
		defer catgut.GetPool().Put(percent)

		percent.Set("%")

		loc, ok := slices.BinarySearchFunc(
			s,
			percent,
			func(ewp EtikettWithParentsAndTypes, e *Etikett) int {
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
		func(ewp EtikettWithParentsAndTypes, e *Etikett) int {
			cmp := catgut.CompareUTF8Bytes(
				ewp.Etikett.Bytes()[offset:],
				e.Bytes(),
				partial,
			)

			return cmp
		},
	)
}

func (s TagsWithParentsAndTypes) ContainsEtikett(e *Etikett) (int, bool) {
	return slices.BinarySearchFunc(
		s,
		e,
		func(ewp EtikettWithParentsAndTypes, e *Etikett) int {
			cmp := ewp.Etikett.ComparePartial(e)
			return cmp
		},
	)
}

func (s TagsWithParentsAndTypes) GetMatching(
	e *Etikett,
) (matching []EtikettWithParentsAndTypes) {
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
func (s *TagsWithParentsAndTypes) Add(
	e1 *Etikett,
	p *PathWithType,
) (err error) {
	var e *Etikett

	if e, err = e1.Clone(); err != nil {
		err = errors.Wrap(err)
		return
	}

	idx, ok := s.ContainsEtikett(e)

	var a EtikettWithParentsAndTypes

	if ok {
		a = (*s)[idx]
		a.Parents.AddNonEmptyPath(p)
		(*s)[idx] = a
	} else {
		a = EtikettWithParentsAndTypes{Etikett: e}
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
func (s *TagsWithParentsAndTypes) Remove(e1 *Etikett) (err error) {
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

func (s TagsWithParentsAndTypes) StringCommaSeparatedExplicit() string {
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
