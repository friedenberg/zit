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
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/kasten"
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
	a.Etiketten = collections.MakeMutableSetStringer[ketikett]()
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
	EtikettenHidden             schnittstellen.Set[kennung.Etikett]
	EtikettenHiddenStringsSlice []string
	EtikettenToAddToNew         []string
	Etiketten                   schnittstellen.MutableSet[ketikett]
	ImplicitEtiketten           implicitEtikettenMap

	// Typen
	ExtensionsToTypen map[string]string
	DefaultTyp        typ.Transacted // deprecated
	Typen             typSet

	// Kasten
	Kisten kastenSet
}

func Make(
	s standort.Standort,
	kcli erworben.Cli,
) (c *Compiled, err error) {
	c = &Compiled{
		cli: kcli,
		compiled: compiled{
			ExtensionsToTypen: make(map[string]string),
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

func (kc *compiled) recompile() (err error) {
	kc.hasChanges = true
	kc.EtikettenToAddToNew = collections.ResetSlice(kc.EtikettenToAddToNew)

	{
		kc.ImplicitEtiketten = make(implicitEtikettenMap)
		etikettenHidden := collections.MakeMutableSetStringer[kennung.Etikett]()

		if err = kc.Etiketten.EachPtr(
			func(ke *ketikett) (err error) {
				ct := ke.Transacted
				k := ct.Sku.GetKennung()
				tn := k.String()
				tv := ct.Akte

				switch {
				case tv.Hide:
					etikettenHidden.Add(k)
					kc.EtikettenHiddenStringsSlice = append(
						kc.EtikettenHiddenStringsSlice,
						tn,
					)

					// TODO-P2: determine why empty etiketten make it here
				case tv.AddToNewZettels && tn != "":
					kc.EtikettenToAddToNew = append(kc.EtikettenToAddToNew, tn)
				}

				// kc.applyExpandedEtikett(ct)

				if err = kc.AccumulateImplicitEtiketten(
					ke.Transacted.Sku.GetKennung(),
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

		sort.Slice(kc.EtikettenToAddToNew, func(i, j int) bool {
			return kc.EtikettenToAddToNew[i] < kc.EtikettenToAddToNew[j]
		})

		kc.EtikettenHidden = etikettenHidden.ImmutableClone()
	}

	if err = kc.Typen.Each(
		func(ct *typ.Transacted) (err error) {
			fe := ct.Akte.FileExtension

			if fe != "" {
				kc.ExtensionsToTypen[fe] = ct.Sku.GetKennung().String()
			}

			if ct == nil {
				errors.Todo("determine why any Typen might be nil")
				return
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

func (kc *Compiled) Flush(s standort.Standort) (err error) {
	if !kc.hasChanges || kc.DryRun {
		return
	}

	if err = kc.recompile(); err != nil {
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
) (expandedActual []*typ.Transacted) {
	expandedMaybe := collections.MakeMutableSetStringer[values.String]()

	sa := collections.MakeFuncSetString[
		values.String,
		*values.String,
	](expandedMaybe)

	typExpander.Expand(sa, v)
	expandedActual = make([]*typ.Transacted, 0)

	expandedMaybe.Each(
		func(v values.String) (err error) {
			c.lock.Lock()
			defer c.lock.Unlock()

			ct, ok := c.Typen.Get(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, ct)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return len(
			expandedActual[i].Sku.GetKennung().String(),
		) > len(
			expandedActual[j].Sku.GetKennung().String(),
		)
	})

	return
}

func (kc compiled) IsInlineTyp(k kennung.Typ) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	tc := kc.GetApproximatedTyp(k)

	if !tc.HasValue() {
		return
	}

	isInline = tc.ApproximatedOrActual().Akte.InlineAkte

	return
}

// Returns the exactly matching Typ, or if it doesn't exist, returns the parent
// Typ or nil. (Parent Typ for `md-gdoc` would be `md`.)
func (kc compiled) GetApproximatedTyp(k kennung.Typ) (ct ApproximatedTyp) {
	expandedActual := kc.GetSortedTypenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct.hasValue = true
		ct.typ = *expandedActual[0]

		if ct.typ.Sku.GetKennung().Equals(k) {
			ct.isActual = true
		}
	}

	return
}

func (kc compiled) GetKasten(k kennung.Kasten) (ct *kasten.Transacted) {
	ct, _ = kc.Kisten.Get(k.String())
	return
}

func (k *compiled) SetTransacted(
	kt *erworben.Transacted,
) {
	if !k.Sku.Less(kt.Sku) {
		return
	}

	k.hasChanges = true
	k.Sku = kt.Sku
	k.Akte = kt.Akte
}

func (k *compiled) AddKasten(
	b *kasten.Transacted,
) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.hasChanges = true
	a, ok := k.Kisten.Get(k.Kisten.Key(b))

	if !ok || a.Less(*b) {
		k.Kisten.Add(b)
	}

	return
}

func (k *compiled) AddTyp(
	b *typ.Transacted,
) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.hasChanges = true

	a, ok := k.Typen.Get(k.Typen.Key(b))

	if !ok || a.Less(*b) {
		k.Typen.Add(b)
	}

	return
}

func (c *compiled) applyExpandedTyp(ct typ.Transacted) {
	expandedActual := c.GetSortedTypenExpanded(ct.Sku.GetKennung().String())

	for _, ex := range expandedActual {
		ct.Akte.Merge(ex.Akte)
	}
}
