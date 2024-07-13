package metadatei

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/flag_policy"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	flag2 "code.linenisgreat.com/zit/go/zit/src/bravo/flag"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type MetadateiWriterTo interface {
	io.WriterTo
	HasMetadateiContent() bool
}

type Metadatei struct {
	// StoreVersion values.Int
	// Domain
	Kasten      kennung.RepoId
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.TagMutableSet // public for gob, but should be private
	Typ         kennung.Type

	Shas
	Tai kennung.Tai

	Comments      []string
	Verzeichnisse Verzeichnisse
}

func (m *Metadatei) GetMetadatei() *Metadatei {
	return m
}

func (m *Metadatei) Sha() *sha.Sha {
	return &m.SelbstMetadateiKennungMutter
}

func (m *Metadatei) Mutter() *sha.Sha {
	return &m.MutterMetadateiKennungMutter
}

func (m *Metadatei) AddToFlagSet(f *flag.FlagSet) {
	f.Var(
		&m.Bezeichnung,
		"bezeichnung",
		"the Bezeichnung to use for created or updated Zettelen",
	)

	// TODO add support for etiketten_path
	fes := flag2.Make(
		flag_policy.FlagPolicyAppend,
		func() string {
			return m.Verzeichnisse.Etiketten.String()
		},
		func(v string) (err error) {
			vs := strings.Split(v, ",")

			for _, v := range vs {
				if err = m.AddEtikettString(v); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
		func() {
			m.ResetEtiketten()
		},
	)

	f.Var(
		fes,
		"etiketten",
		"the Etiketten to use for created or updated Objekte",
	)

	f.Func(
		"typ",
		"the Typ for the created or updated Objekte",
		func(v string) (err error) {
			return m.Typ.Set(v)
		},
	)
}

func (z *Metadatei) UserInputIsEmpty() bool {
	if !z.Bezeichnung.IsEmpty() {
		return false
	}

	if z.Etiketten != nil && z.Etiketten.Len() > 0 {
		return false
	}

	if !kennung.IsEmpty(z.Typ) {
		return false
	}

	return true
}

func (z *Metadatei) IsEmpty() bool {
	if !z.Akte.IsNull() {
		return false
	}

	if !z.UserInputIsEmpty() {
		return false
	}

	if !z.Tai.IsZero() {
		return false
	}

	return true
}

func (z *Metadatei) SetBezeichnung(b bezeichnung.Bezeichnung) {
	z.Bezeichnung = b
}

func (z *Metadatei) SetTyp(t kennung.Type) {
	z.Typ = t
}

func (z *Metadatei) GetBezeichnung() bezeichnung.Bezeichnung {
	return z.Bezeichnung
}

func (z *Metadatei) GetBezeichnungPtr() *bezeichnung.Bezeichnung {
	return &z.Bezeichnung
}

func (m *Metadatei) GetEtiketten() kennung.TagSet {
	if m.Etiketten == nil {
		m.Etiketten = kennung.MakeTagMutableSet()
	}

	return m.Etiketten
}

func (m *Metadatei) ResetEtiketten() {
	if m.Etiketten == nil {
		m.Etiketten = kennung.MakeTagMutableSet()
	}

	m.Etiketten.Reset()
	m.Verzeichnisse.Etiketten.Reset()
}

func (z *Metadatei) AddEtikettString(es string) (err error) {
	if es == "" {
		return
	}

	var e kennung.Tag

	if err = e.Set(es); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = z.AddEtikettPtr(&e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m *Metadatei) AddEtikettPtr(e *kennung.Tag) (err error) {
	if e == nil || e.String() == "" {
		return
	}

	if m.Etiketten == nil {
		m.Etiketten = kennung.MakeTagMutableSet()
	}

	kennung.AddNormalizedEtikett(m.Etiketten, e)
	cs := catgut.MakeFromString(e.String())
	m.Verzeichnisse.Etiketten.AddEtikett(cs)

	return
}

func (m *Metadatei) AddEtikettPtrFast(e *kennung.Tag) (err error) {
	if m.Etiketten == nil {
		m.Etiketten = kennung.MakeTagMutableSet()
	}

	if err = m.Etiketten.Add(*e); err != nil {
		err = errors.Wrap(err)
		return
	}

	cs := catgut.MakeFromString(e.String())

	if err = m.Verzeichnisse.Etiketten.AddEtikett(cs); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m *Metadatei) SetEtiketten(e kennung.TagSet) {
	if m.Etiketten == nil {
		m.Etiketten = kennung.MakeTagMutableSet()
	}

	iter.ResetMutableSetWithPool(m.Etiketten, kennung.GetTagPool())

	if e == nil {
		return
	}

	if e.Len() == 1 && e.Any().String() == "" {
		panic("empty etikett set")
	}

	errors.PanicIfError(e.EachPtr(m.AddEtikettPtr))
}

func (z *Metadatei) SetAkteSha(sh interfaces.ShaGetter) {
	z.Akte.SetShaLike(sh)
}

func (z *Metadatei) GetTyp() kennung.Type {
	return z.Typ
}

func (z *Metadatei) GetTypPtr() *kennung.Type {
	return &z.Typ
}

func (z *Metadatei) GetTai() kennung.Tai {
	return z.Tai
}

// TODO-P2 remove
func (b *Metadatei) EqualsSansTai(a *Metadatei) bool {
	return EqualerSansTai.Equals(a, b)
}

// TODO-P2 remove
func (pz *Metadatei) Equals(z1 *Metadatei) bool {
	return Equaler.Equals(pz, z1)
}

func (a *Metadatei) Subtract(
	b *Metadatei,
) {
	if a.Typ.String() == b.Typ.String() {
		a.Typ = kennung.Type{}
	}

	err := b.GetEtiketten().EachPtr(
		func(e *kennung.Tag) (err error) {
			return a.Etiketten.DelPtr(e)
		},
	)
	errors.PanicIfError(err)
}

func (mp *Metadatei) AddComment(f string, vals ...interface{}) {
	mp.Comments = append(mp.Comments, fmt.Sprintf(f, vals...))
}

func (selbst *Metadatei) SetMutter(mg Getter) (err error) {
	mutter := mg.GetMetadatei()

	if err = selbst.Mutter().SetShaLike(
		mutter.Sha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = selbst.MutterMetadateiKennungMutter.SetShaLike(
		&mutter.SelbstMetadateiKennungMutter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m *Metadatei) GenerateExpandedEtiketten() {
	m.Verzeichnisse.SetExpandedEtiketten(kennung.ExpandMany(
		m.GetEtiketten(),
		expansion.ExpanderRight,
	))
}
