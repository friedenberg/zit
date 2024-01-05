package etiketten_path

import (
	"bytes"
	"io"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/delta/ohio"
)

type etikettLike interface {
	schnittstellen.Stringer
	Bytes() []byte
}

type Path []*catgut.String

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

func (p *Path) String() string {
	var sb strings.Builder

	afterFirst := false
	for i := p.Len() - 1; i >= 0; i-- {
		if afterFirst {
			sb.WriteByte(' ')
			sb.WriteByte('-')
			sb.WriteByte('>')
			sb.WriteByte(' ')
		}

		afterFirst = true

		s := (*p)[i]
		sb.Write(s.Bytes())
	}

	return sb.String()
}

func (a *Path) Copy() (b *Path) {
	b = &Path{}
	*b = make([]*catgut.String, a.Len())

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
	var x catgut.String
	x.SetBytes(a.Bytes())
	a.SetBytes(b.Bytes())
	b.SetBytes(x.Bytes())
}

func (p *Path) Add(e etikettLike) {
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
		*p = append(*p, make([]*catgut.String, diff)...)
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
