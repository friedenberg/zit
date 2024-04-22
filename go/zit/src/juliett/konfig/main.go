package konfig

import (
	"encoding/gob"
	"fmt"
	"sort"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	pkg_angeboren "code.linenisgreat.com/zit/src/delta/angeboren"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/akten"
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

func (a *compiled) Reset() error {
	a.ExtensionsToTypen = make(map[string]string)
	a.TypenToExtensions = make(map[string]string)

	a.lock = &sync.Mutex{}
	a.EtikettenHidden = kennung.MakeEtikettSet()
	a.Etiketten = collections_value.MakeMutableValueSet[*ketikett](nil)
	a.InlineTypen = collections_value.MakeMutableValueSet[values.String](
		nil,
	)
	a.ImplicitEtiketten = make(implicitEtikettenMap)
	a.Kisten = sku.MakeTransactedMutableSet()
	a.Typen = sku.MakeTransactedMutableSet()

	sku.TransactedResetter.Reset(&a.Sku)

	return nil
}

func (a *Compiled) GetErworben() *erworben.Akte {
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
	Etiketten                   schnittstellen.MutableSetLike[*ketikett]
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

func (c *Compiled) Initialize(
	s standort.Standort,
	kcli erworben.Cli,
) (err error) {
	c.cli = kcli
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

func (kc *Compiled) GetAngeboren() schnittstellen.Angeboren {
	return kc.angeboren
}

func (kc *Compiled) Cli() erworben.Cli {
	return kc.cli
}

func (kc *Compiled) SetCli(k erworben.Cli) {
	kc.cli = k
}

func (kc *Compiled) SetCliFromCommander(k erworben.Cli) {
	oldBasePath := kc.BasePath
	kc.cli = k
	kc.BasePath = oldBasePath
}

func (c *compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.FileExtensions.Zettel)
}

// TODO-P3 merge all the below
func (c *compiled) GetSortedTypenExpanded(
	v string,
) (expandedActual []*sku.Transacted) {
	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)

	sa := iter.MakeFuncSetString(expandedMaybe)

	typExpander.Expand(sa, v)
	expandedActual = make([]*sku.Transacted, 0)

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
			expandedActual[i].GetKennung().String(),
		) > len(
			expandedActual[j].GetKennung().String(),
		)
	})

	return
}

func (kc *compiled) IsInlineTyp(k kennung.Typ) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	isInline = kc.InlineTypen.ContainsKey(k.String())

	return
}

type ApproximatedTyp = akten.ApproximatedTyp

// Returns the exactly matching Typ, or if it doesn't exist, returns the parent
// Typ or nil. (Parent Typ for `md-gdoc` would be `md`.)
func (kc *compiled) GetApproximatedTyp(
	k kennung.Kennung,
) (ct ApproximatedTyp) {
	expandedActual := kc.GetSortedTypenExpanded(k.String())
	if len(expandedActual) > 0 {
		ct.HasValue = true
		ct.Typ = expandedActual[0]

		if kennung.Equals(ct.Typ.GetKennung(), k) {
			ct.IsActual = true
		}
	}

	return
}

func (kc *compiled) GetKasten(k kennung.Kasten) (ct *sku.Transacted) {
	if ct1, ok := kc.Kisten.Get(k.String()); ok {
		ct = sku.GetTransactedPool().Get()
		errors.PanicIfError(ct.SetFromSkuLike(ct1))
	}

	return
}

func (k *compiled) SetTransacted(
	kt1 *sku.Transacted,
	kag schnittstellen.AkteGetter[*erworben.Akte],
) (err error) {
	if !sku.TransactedLessor.LessPtr(&k.Sku, kt1) {
		return
	}

	k.hasChanges = true

	if err = k.Sku.SetFromSkuLike(kt1); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	if err = b.SetFromSkuLike(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	_, err = iter.AddOrReplaceIfGreater(
		k.Kisten,
		b,
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *compiled) AddTransacted(
	a *sku.Transacted,
	ak *akten.Akten,
) (err error) {
	switch a.Kennung.GetGattung() {
	case gattung.Typ:
		return k.AddTyp(a)

	case gattung.Etikett:
		return k.AddEtikett(a)

	case gattung.Kasten:
		return k.AddKasten(a)

	case gattung.Konfig:
		return k.SetTransacted(
			a,
			ak.GetKonfigV0(),
		)
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

	if err = b.SetFromSkuLike(b1); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.lock.Lock()
	defer k.lock.Unlock()

	shouldAdd, err := iter.AddOrReplaceIfGreater(
		k.Typen,
		b,
	)

	k.hasChanges = k.hasChanges || shouldAdd

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
