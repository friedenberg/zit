package metadatei

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type MetadateiWriterTo interface {
	io.WriterTo
	HasMetadateiContent() bool
}

type Metadatei struct {
	// StoreVersion values.Int
	// Kasten
	// Domain
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettMutableSet // public for gob, but should be private
	Typ         kennung.Typ
	Tai         kennung.Tai

	Shas

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

	fes := collections_ptr.MakeFlagCommasFromExisting[kennung.Etikett](
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

func (z *Metadatei) AddEtikettPtr(e *kennung.Etikett) (err error) {
	return iter.AddClonePool[kennung.Etikett, *kennung.Etikett](
		z.GetEtikettenMutable(),
		kennung.GetEtikettPool(),
		kennung.EtikettResetter,
		e,
	)
}

func (z *Metadatei) SetEtiketten(e kennung.EtikettSet) {
	es := z.GetEtikettenMutable()
	iter.ResetMutableSetWithPool(es, kennung.GetEtikettPool())

	if e == nil {
		return
	}

	errors.PanicIfError(
		e.EachPtr(
			iter.MakeAddClonePoolFunc[kennung.Etikett, *kennung.Etikett](
				es,
				kennung.GetEtikettPool(),
				kennung.EtikettResetter,
			),
		),
	)
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
func (pz *Metadatei) EqualsSansTai(z1 *Metadatei) bool {
	return EqualerSansTai.Equals(pz, z1)
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
		d = iter.StringCommaSeparated[kennung.Etikett](z.GetEtiketten())
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

func (z *Metadatei) ApplyGoldenChild(
	e kennung.Etikett,
	mode etikett_rule.RuleGoldenChild,
) (err error) {
	if z.GetEtiketten().Len() == 0 {
		return
	}

	switch mode {
	case etikett_rule.RuleGoldenChildUnset:
		return
	}

	mes := z.GetEtikettenMutable()

	prefixes := iter.Elements[kennung.Etikett](kennung.Withdraw(mes, e))

	if len(prefixes) == 0 {
		return
	}

	var sortFunc func(i, j int) bool

	switch mode {
	case etikett_rule.RuleGoldenChildLowest:
		sortFunc = func(i, j int) bool { return kennung.Less(prefixes[j], prefixes[i]) }

	case etikett_rule.RuleGoldenChildHighest:
		sortFunc = func(i, j int) bool { return kennung.Less(prefixes[i], prefixes[j]) }
	}

	sort.Slice(prefixes, sortFunc)

	mes.Add(prefixes[0])

	return
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
