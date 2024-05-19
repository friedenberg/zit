package etiketten_path

import (
	"bytes"
	"io"
	"sort"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/src/delta/catgut"
)

type Path []*Etikett

func MakePath(els ...*Etikett) *Path {
	p := Path(make([]*Etikett, 0, len(els)))

	for _, e := range els {
		p.Add(e)
	}

	return &p
}

func (a *Path) IsEmpty() bool {
	if a == nil {
		return true
	}

	return a.Len() == 0
}

func (a *Path) First() *Etikett {
	return (*a)[0]
}

func (a *Path) Last() *Etikett {
	return (*a)[a.Len()-1]
}

func (a *Path) Equals(b *Path) bool {
	if a.Len() != b.Len() {
		return false
	}

	for i, as := range *a {
		if !as.Equals((*b)[i]) {
			return false
		}
	}

	return true
}

func (a *Path) Compare(b *Path) int {
	elsA := *a
	elsB := *b

	for {
		lenA, lenB := len(elsA), len(elsB)

		switch {
		case lenA == 0 && lenB == 0:
			return 0

		case lenA == 0:
			return -1

		case lenB == 0:
			return 1
		}

		elA := elsA[0]
		elsA = elsA[1:]

		elB := elsB[0]
		elsB = elsB[1:]

		cmp := elA.Compare(elB)

		if cmp != 0 {
			return cmp
		}
	}

	return 0
}

func (p *Path) String() string {
	var sb strings.Builder

	afterFirst := false
	for i := p.Len() - 1; i >= 0; i-- {
		if afterFirst {
			sb.WriteString(" -> ")
		}

		afterFirst = true

		s := (*p)[i]
		sb.Write(s.Bytes())
	}

	return sb.String()
}

func (a *Path) Copy() (b *Path) {
	b = &Path{}
	*b = make([]*Etikett, a.Len())

	if a == nil {
		return
	}

	for i, s := range *a {
		sb := catgut.GetPool().Get()
		s.CopyTo(sb)
		(*b)[i] = sb
	}

	return
}

func (p *Path) Len() int {
	if p == nil {
		return 0
	}

	return len(*p)
}

func (p *Path) Cap() int {
	if p == nil {
		return 0
	}

	return cap(*p)
}

func (p *Path) Less(i, j int) bool {
	return bytes.Compare((*p)[i].Bytes(), (*p)[i].Bytes()) == -1
}

func (p *Path) Swap(i, j int) {
	a, b := (*p)[i], (*p)[j]
	var x Etikett
	x.SetBytes(a.Bytes())
	a.SetBytes(b.Bytes())
	b.SetBytes(x.Bytes())
}

func (p *Path) Add(e *Etikett) {
	if e.IsEmpty() {
		return
	}

	*p = append(*p, catgut.GetPool().Get())
	(*p)[p.Len()-1].SetBytes(e.Bytes())
	sort.Sort(p)
}

func (p *Path) ReadFrom(r io.Reader) (n int64, err error) {
	var count uint8

	var n1 int
	if count, n1, err = ohio.ReadUint8(r); err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n += int64(n1)

	*p = (*p)[:p.Cap()]

	if diff := count - uint8(p.Len()); diff > 0 {
		*p = append(*p, make([]*Etikett, diff)...)
	}

	for i := uint8(0); i < count; i++ {
		var cl uint8

		if cl, n1, err = ohio.ReadUint8(r); err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}

		n += int64(n1)

		if (*p)[i] == nil {
			(*p)[i] = catgut.GetPool().Get()
		}

		_, err = (*p)[i].ReadNFrom(r, int(cl))
		if err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}
	}

	return
}

func (p *Path) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int

	n1, err = ohio.WriteUint8(w, uint8(p.Len()))
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	for _, s := range *p {
		if s.Len() == 0 {
			panic("found empty etikett in etiketten_path")
		}

		n1, err = ohio.WriteUint8(w, uint8(s.Len()))
		n += int64(n1)

		if err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}

		var n2 int64
		n2, err = s.WriteTo(w)
		n += n2

		if err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}
	}

	return
}
