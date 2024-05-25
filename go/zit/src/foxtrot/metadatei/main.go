package metadatei

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type MetadateiWriterTo interface {
	io.WriterTo
	HasMetadateiContent() bool
}

type Metadatei struct {
	// StoreVersion values.Int
	// Domain
	Kasten      kennung.Kasten
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettMutableSet // public for gob, but should be private
	Typ         kennung.Typ

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
	fes := collections_ptr.MakeFlagCommasFromExisting(
		collections_ptr.SetterPolicyAppend,
		m.GetEtikettenMutable(),
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

func (z *Metadatei) SetTyp(t kennung.Typ) {
	z.Typ = t
}

func (z *Metadatei) GetBezeichnung() bezeichnung.Bezeichnung {
	return z.Bezeichnung
}

func (z *Metadatei) GetBezeichnungPtr() *bezeichnung.Bezeichnung {
	return &z.Bezeichnung
}

func (z *Metadatei) GetEtiketten() kennung.EtikettSet {
	return z.GetEtikettenMutable()
}

func (z *Metadatei) GetEtikettenMutable() kennung.EtikettMutableSet {
	if z.Etiketten == nil {
		z.Etiketten = kennung.MakeEtikettMutableSet()
	}

	return z.Etiketten
}

func (m *Metadatei) ResetEtiketten() {
	m.Etiketten.Reset()
	m.Verzeichnisse.Etiketten.Reset()
}

func (z *Metadatei) AddEtikettString(es string) (err error) {
	var e kennung.Etikett

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

func (z *Metadatei) AddEtikettPtr(e *kennung.Etikett) (err error) {
	kennung.AddNormalizedEtikett(z.GetEtikettenMutable(), e)
	z.Verzeichnisse.Etiketten.AddEtikett(catgut.MakeFromString(e.String()))
	return
}

func (m *Metadatei) SetEtiketten(e kennung.EtikettSet) {
	if m.Etiketten == nil {
		m.Etiketten = kennung.MakeEtikettMutableSet()
	}

	iter.ResetMutableSetWithPool(m.Etiketten, kennung.GetEtikettPool())

	if e == nil {
		return
	}

	errors.PanicIfError(e.EachPtr(m.AddEtikettPtr))
}

func (z *Metadatei) SetAkteSha(sh schnittstellen.ShaGetter) {
	z.Akte.SetShaLike(sh)
}

func (z *Metadatei) GetTyp() kennung.Typ {
	return z.Typ
}

func (z *Metadatei) GetTypPtr() *kennung.Typ {
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

func (z *Metadatei) String() (d string) {
	return z.Description()
}

func (z *Metadatei) Description() (d string) {
	d = z.Bezeichnung.String()

	if strings.TrimSpace(d) == "" {
		d = iter.StringCommaSeparated(z.GetEtiketten())
	}

	return
}

func (a *Metadatei) Subtract(
	b *Metadatei,
) {
	if a.Typ.String() == b.Typ.String() {
		a.Typ = kennung.Typ{}
	}

	err := b.GetEtiketten().EachPtr(
		func(e *kennung.Etikett) (err error) {
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
