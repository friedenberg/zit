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
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/iter2"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/srx/bravo/expansion"
)

var typExpander expansion.Expander

func init() {
	typExpander = expansion.MakeExpanderRight(`-`)

	gob.Register(
		collections_value.MakeMutableValueSet[values.String](
			nil,
		),
	)

	gob.Register(
		collections_value.MakeValueSet[values.String](
			nil,
		),
	)

	gob.Register(iter.StringerKeyer[values.String]{})
	gob.Register(iter.StringerKeyerPtr[kennung.Typ, *kennung.Typ]{})
}

type angeboren = pkg_angeboren.Konfig

type Compiled struct {
	cli
	compiled
	angeboren
}

func (a *compiled) Reset() {
	a.lock = &sync.Mutex{}
	a.EtikettenHidden = kennung.MakeEtikettSet()
	a.Etiketten = collections_ptr.MakeMutableValueSet[ketikett, *ketikett](nil)
	a.InlineTypen = collections_value.MakeMutableValueSet[values.String](
		nil,
	)
	a.ImplicitEtiketten = make(implicitEtikettenMap)
	a.Kisten = sku.MakeTransactedMutableSet()
	a.Typen = sku.MakeTransactedMutableSet()
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

	Sku sku.Transacted

	erworben.Akte

	// Etiketten
	EtikettenHidden             kennung.EtikettSet
	EtikettenHiddenStringsSlice []string
	DefaultEtiketten            kennung.EtikettSet
	Etiketten                   schnittstellen.MutableSetPtrLike[ketikett, *ketikett]
	ImplicitEtiketten           implicitEtikettenMap

	// Typen
	ExtensionsToTypen map[string]string
	TypenToExtensions map[string]string
	DefaultTyp        sku.Transacted // deprecated
	Typen             sku.TransactedMutableSet
	InlineTypen       schnittstellen.SetLike[values.String]

	// Kasten
	Kisten sku.TransactedMutableSet
}

func Make(
	s standort.Standort,
	kcli erworben.Cli,
) (c *Compiled, err error) {
	c = &Compiled{
		cli: kcli,
		compiled: compiled{
			ExtensionsToTypen: make(map[string]string),
			TypenToExtensions: make(map[string]string),
		},
	}

	c.Reset()

	c.angeboren = s.GetKonfig()

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

func (kc Compiled) HasChanges() bool {
	kc.lock.Lock()
	defer kc.lock.Unlock()

	return kc.hasChanges
}

func (kc *Compiled) SetHasChanges(v bool) {
	kc.lock.Lock()
	defer kc.lock.Unlock()

	kc.hasChanges = true
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
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	kc.DefaultEtiketten = kennung.MakeEtikettSet(kc.Akte.Defaults.Etiketten...)

	{
		kc.ImplicitEtiketten = make(implicitEtikettenMap)

		if err = kc.Etiketten.Each(
			func(ke ketikett) (err error) {
				var e kennung.Etikett

				if err = e.Set(ke.Transacted.GetKennung().String()); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = kc.AccumulateImplicitEtiketten(e); err != nil {
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

	inlineTypen := collections_value.MakeMutableValueSet[values.String](nil)

	defer func() {
		kc.InlineTypen = inlineTypen.CloneSetLike()
	}()

	if err = kc.Typen.EachPtr(
		func(ct *sku.Transacted) (err error) {
			var ta *typ_akte.V0

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
			kc.ExtensionsToTypen[fe] = ct.GetKennung().String()
			kc.TypenToExtensions[ct.GetKennung().String()] = fe

			if ta.InlineAkte {
				inlineTypen.Add(values.MakeString(ct.Kennung.String()))
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
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
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
) (expandedActual []*sku.Transacted) {
	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)

	sa := iter.MakeFuncSetString[
		values.String,
		*values.String,
	](expandedMaybe)

	typExpander.Expand(sa, v)
	expandedActual = make([]*sku.Transacted, 0)

	expandedMaybe.Each(
		func(v values.String) (err error) {
			c.lock.Lock()
			defer c.lock.Unlock()

			ct, ok := c.Typen.GetPtr(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, ct)

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

	isInline = kc.InlineTypen.ContainsKey(k.String())

	return
}

// Returns the exactly matching Typ, or if it doesn't exist, returns the parent
// Typ or nil. (Parent Typ for `md-gdoc` would be `md`.)
func (kc compiled) GetApproximatedTyp(
	k kennung.Kennung,
) (ct ApproximatedTyp) {
	expandedActual := kc.GetSortedTypenExpanded(k.String())
	if len(expandedActual) > 0 {
		ct.hasValue = true
		ct.typ = expandedActual[0]

		if kennung.Equals(ct.typ.GetKennung(), k) {
			ct.isActual = true
		}
	}

	return
}

func (kc compiled) GetKasten(k kennung.Kasten) (ct *sku.Transacted) {
	if ct1, ok := kc.Kisten.GetPtr(k.String()); ok {
		ct = sku.GetTransactedPool().Get()
		errors.PanicIfError(ct.SetFromSkuLike(ct1))
	}

	return
}

func (k *compiled) SetTransacted(
	kt1 *sku.Transacted,
	kag schnittstellen.AkteGetter[*erworben.Akte],
) (err error) {
	kt := sku.GetTransactedPool().Get()
	*kt = *kt1

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
	c *sku.Transacted,
) (err error) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.hasChanges = true

	b := sku.GetTransactedPool().Get()
	*b = *c

	err = iter2.AddPtrOrReplaceIfGreater[sku.Transacted, *sku.Transacted](
		k.Kisten,
		sku.TransactedLessor,
		b,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *compiled) AddTyp(
	b1 *sku.Transacted,
) (err error) {
	if err = gattung.Typ.AssertGattung(b1); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := sku.GetTransactedPool().Get()
	*b = *b1

	k.lock.Lock()
	defer k.lock.Unlock()

	k.hasChanges = true

	err = iter2.AddPtrOrReplaceIfGreater[sku.Transacted, *sku.Transacted](
		k.Typen,
		sku.TransactedLessor,
		b,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *compiled) applyExpandedTyp(ct sku.Transacted) {
}