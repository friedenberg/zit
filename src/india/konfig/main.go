package konfig

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"sync"

	pkg_angeboren "github.com/friedenberg/zit/src/alfa/angeboren"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/typ"
)

var typExpander kennung.Expander

func init() {
	typExpander = kennung.MakeExpanderRight(`-`)
}

type angeboren = pkg_angeboren.Konfig

type Compiled struct {
	cli
	compiled
	angeboren
}

func (a *compiled) Reset() {
	a.lock = &sync.Mutex{}
	a.Typen = makeCompiledTypSetFromSlice(nil)
	a.EtikettenHidden = kennung.MakeEtikettSet()
	a.Etiketten = collections_ptr.MakeMutableValueSet[ketikett, *ketikett](nil)
	a.InlineTypen = collections_ptr.MakeMutableValueSet[kennung.Typ, *kennung.Typ](
		nil,
	)
	a.ImplicitEtiketten = make(implicitEtikettenMap)
	a.Kisten = makeCompiledKastenSetFromSlice(nil)
}

func (a Compiled) GetErworben() erworben.Akte {
	return a.Akte
}

func (a *Compiled) GetErworbenPtr() *erworben.Akte {
	return &a.Akte
}

type cli = erworben.Cli

type compiled struct {
	lock sync.Locker

	hasChanges bool

	Sku sku.Transacted[kennung.Konfig, *kennung.Konfig]

	erworben.Akte

	// Etiketten
	EtikettenHidden             kennung.EtikettSet
	EtikettenHiddenStringsSlice []string
	DefaultEtiketten            kennung.EtikettSet
	Etiketten                   schnittstellen.MutableSetPtrLike[ketikett, *ketikett]
	ImplicitEtiketten           implicitEtikettenMap

	// Typen
	ExtensionsToTypen map[string]kennung.Typ
	TypenToExtensions map[kennung.Typ]string
	DefaultTyp        sku.TransactedTyp // deprecated
	Typen             schnittstellen.MutableSetLike[sku.TransactedTyp]
	InlineTypen       schnittstellen.SetPtrLike[kennung.Typ, *kennung.Typ]

	// Kasten
	Kisten schnittstellen.MutableSetLike[sku.TransactedKasten]
}

func Make(
	s standort.Standort,
	kcli erworben.Cli,
) (c *Compiled, err error) {
	c = &Compiled{
		cli: kcli,
		compiled: compiled{
			ExtensionsToTypen: make(map[string]kennung.Typ),
			TypenToExtensions: make(map[kennung.Typ]string),
		},
	}

	c.Reset()

	if err = c.loadKonfigAngeboren(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.loadKonfigErworben(s); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (kc *Compiled) loadKonfigAngeboren(s standort.Standort) (err error) {
	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(s.FileKonfigAngeboren()); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewDecoder(f)

	if err = dec.Decode(&kc.angeboren); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *Compiled) loadKonfigErworben(s standort.Standort) (err error) {
	var f *os.File

	p := s.FileKonfigCompiled()

	if kc.angeboren.UseKonfigErworbenFile {
		p = s.FileKonfigErworben()
	}

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewDecoder(f)

	if err = dec.Decode(&kc.compiled); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (kc Compiled) GetAngeboren() schnittstellen.Angeboren {
	return kc.angeboren
}

func (kc Compiled) Cli() erworben.Cli {
	return kc.cli
}

func (kc *Compiled) SetCli(k erworben.Cli) {
	kc.cli = k
}

func (kc *Compiled) SetCliFromCommander(k erworben.Cli) {
	oldBasePath := kc.cli.BasePath
	kc.cli = k
	kc.cli.BasePath = oldBasePath
}

func (kc *compiled) recompile(
	tagp schnittstellen.AkteGetterPutter[*typ.Akte],
) (err error) {
	kc.DefaultEtiketten = kennung.MakeEtikettSet(kc.Akte.Defaults.Etiketten...)

	{
		kc.ImplicitEtiketten = make(implicitEtikettenMap)

		if err = kc.Etiketten.Each(
			func(ke ketikett) (err error) {
				if err = kc.AccumulateImplicitEtiketten(
					ke.Transacted.GetKennung(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		sort.Slice(kc.EtikettenHiddenStringsSlice, func(i, j int) bool {
			return kc.EtikettenHiddenStringsSlice[i] < kc.EtikettenHiddenStringsSlice[j]
		})
	}

	{
		kc.EtikettenHidden = kennung.MakeEtikettSet(
			kc.Akte.HiddenEtiketten...,
		)
	}

	inlineTypen := kc.InlineTypen.CloneMutableSetPtrLike()

	defer func() {
		kc.InlineTypen = inlineTypen.CloneSetPtrLike()
	}()

	if err = kc.Typen.Each(
		func(ct sku.TransactedTyp) (err error) {
			var ta *typ.Akte

			if ta, err = tagp.GetAkte(ct.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer tagp.PutAkte(ta)

			fe := ta.FileExtension

			if fe == "" {
				fe = ct.GetKennung().String()
			}

			// TODO-P2 enforce uniqueness
			kc.ExtensionsToTypen[fe] = ct.GetKennung()
			kc.TypenToExtensions[ct.GetKennung()] = fe

			if ta.InlineAkte {
				inlineTypen.Add(ct.Kennung)
			}

			// kc.applyExpandedTyp(*ct)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *Compiled) Flush(
	s standort.Standort,
	tagp schnittstellen.AkteGetterPutter[*typ.Akte],
) (err error) {
	if !kc.hasChanges || kc.DryRun {
		return
	}

	if err = kc.recompile(tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := s.FileKonfigCompiled()

	if kc.angeboren.UseKonfigErworbenFile {
		p = s.FileKonfigErworben()
	}

	var f *os.File

	if f, err = files.OpenExclusiveWriteOnlyTruncate(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	dec := gob.NewEncoder(f)

	if err = dec.Encode(kc.compiled); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.FileExtensions.Zettel)
}

// TODO-P3 merge all the below
func (c compiled) GetSortedTypenExpanded(
	v string,
) (expandedActual []*sku.TransactedTyp) {
	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)

	sa := iter.MakeFuncSetString[
		values.String,
		*values.String,
	](expandedMaybe)

	typExpander.Expand(sa, v)
	expandedActual = make([]*sku.TransactedTyp, 0)

	expandedMaybe.Each(
		func(v values.String) (err error) {
			c.lock.Lock()
			defer c.lock.Unlock()

			ct, ok := c.Typen.Get(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, &ct)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return len(
			expandedActual[i].GetKennung().String(),
		) > len(
			expandedActual[j].GetKennung().String(),
		)
	})

	return
}

func (kc compiled) IsInlineTyp(k kennung.Typ) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	isInline = kc.InlineTypen.Contains(k)

	return
}

// Returns the exactly matching Typ, or if it doesn't exist, returns the parent
// Typ or nil. (Parent Typ for `md-gdoc` would be `md`.)
func (kc compiled) GetApproximatedTyp(k kennung.Typ) (ct ApproximatedTyp) {
	expandedActual := kc.GetSortedTypenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct.hasValue = true
		ct.typ = *expandedActual[0]

		if ct.typ.GetKennung().Equals(k) {
			ct.isActual = true
		}
	}

	return
}

func (kc compiled) GetKasten(k kennung.Kasten) (ct *sku.TransactedKasten) {
	var ct1 sku.TransactedKasten
	ct1, _ = kc.Kisten.Get(k.String())
	ct = &ct1
	return
}

func (k *compiled) SetTransacted(
	kt *erworben.Transacted,
	kag schnittstellen.AkteGetter[*erworben.Akte],
) (err error) {
	if !k.Sku.Less(*kt) {
		return
	}

	k.hasChanges = true
	k.Sku = *kt

	var a *erworben.Akte

	if a, err = kag.GetAkte(k.Sku.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.Akte = *a

	return
}

func (k *compiled) AddKasten(
	b *sku.TransactedKasten,
) (err error) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.hasChanges = true

	if err = iter.AddOrReplaceIfGreater(k.Kisten, *b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *compiled) AddTyp(
	b *sku.TransactedTyp,
) (err error) {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.hasChanges = true

	if err = iter.AddOrReplaceIfGreater(k.Typen, *b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *compiled) applyExpandedTyp(ct sku.TransactedTyp) {
}
