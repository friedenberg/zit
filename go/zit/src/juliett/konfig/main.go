package konfig

import (
	"encoding/gob"
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
	"code.linenisgreat.com/zit/src/juliett/query"
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
	schlummernd *query.Schlummernd
}

func (a *compiled) Reset() error {
	a.ExtensionsToTypen = make(map[string]string)
	a.TypenToExtensions = make(map[string]string)

	a.lock = &sync.Mutex{}
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
	DefaultEtiketten  kennung.EtikettSet
	Etiketten         schnittstellen.MutableSetLike[*ketikett]
	ImplicitEtiketten implicitEtikettenMap

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
  schlummernd *query.Schlummernd,
) (err error) {
	c.cli = kcli
	c.Reset()
	c.angeboren = s.GetKonfig()
  c.schlummernd = schlummernd

	wg := iter.MakeErrorWaitGroupParallel()
	wg.Do(func() (err error) {
		if err = c.loadKonfigErworben(s); err != nil {
			if errors.IsNotExist(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}
		}

		return
	})

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *Compiled) SetCli(k erworben.Cli) {
	kc.cli = k
}

func (kc *Compiled) SetCliFromCommander(k erworben.Cli) {
	oldBasePath := kc.BasePath
	kc.cli = k
	kc.BasePath = oldBasePath
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

func (k *compiled) SetTransacted(
	kt1 *sku.Transacted,
	kag schnittstellen.AkteGetter[*erworben.Akte],
) (err error) {
	if !sku.TransactedLessor.LessPtr(&k.Sku, kt1) {
		return
	}

	k.lock.Lock()
	defer k.lock.Unlock()

	k.setHasChanges()

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
	k.setHasChanges()

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

func (k *Compiled) ApplyAndAddTransacted(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	ak *akten.Akten,
) (err error) {
	if err = k.ApplyToSku(kinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch kinder.Kennung.GetGattung() {
	case gattung.Typ:
		return k.AddTyp(kinder)

	case gattung.Etikett:
		return k.compiled.AddEtikett(kinder, mutter)

	case gattung.Kasten:
		return k.AddKasten(kinder)

	case gattung.Konfig:
		return k.SetTransacted(
			kinder,
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

	if shouldAdd {
		k.setHasChanges()
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
