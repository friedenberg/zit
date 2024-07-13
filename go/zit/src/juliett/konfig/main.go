package konfig

import (
	"encoding/gob"
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	pkg_angeboren "code.linenisgreat.com/zit/go/zit/src/delta/angeboren"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/akten"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
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

	changes []string

	Sku sku.Transacted

	erworben.Akte

	// Etiketten
	DefaultEtiketten  kennung.EtikettSet
	Etiketten         interfaces.MutableSetLike[*ketikett]
	ImplicitEtiketten implicitEtikettenMap

	// Typen
	ExtensionsToTypen map[string]string
	TypenToExtensions map[string]string
	DefaultTyp        sku.Transacted // deprecated
	Typen             sku.TransactedMutableSet
	InlineTypen       interfaces.SetLike[values.String]

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

func (k *compiled) setTransacted(
	kt1 *sku.Transacted,
	kag interfaces.BlobGetter[*erworben.Akte],
) (didChange bool, err error) {
	if !sku.TransactedLessor.LessPtr(&k.Sku, kt1) {
		return
	}

	k.lock.Lock()
	defer k.lock.Unlock()

	didChange = true

	if err = k.Sku.SetFromSkuLike(kt1); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.setHasChanges(fmt.Sprintf("updated konfig: %s", &k.Sku))

	var a *erworben.Akte

	if a, err = kag.GetBlob(k.Sku.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.Akte = *a

	return
}

func (k *compiled) addKasten(
	c *sku.Transacted,
) (didChange bool, err error) {
	k.lock.Lock()
	defer k.lock.Unlock()

	b := sku.GetTransactedPool().Get()

	if err = b.SetFromSkuLike(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	if didChange, err = iter.AddOrReplaceIfGreater(
		k.Kisten,
		b,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *Compiled) IsDryRun() bool {
	return k.DryRun
}

func (k *Compiled) GetTypStringFromExtension(t string) string {
	return k.ExtensionsToTypen[t]
}

func (k *Compiled) GetTypExtension(v string) string {
	return k.TypenToExtensions[v]
}

func (k *Compiled) AddTransacted(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	ak *akten.Akten,
	mode objekte_mode.Mode,
) (err error) {
	didChange := false

	switch kinder.Kennung.GetGenre() {
	case gattung.Typ:
		if didChange, err = k.addTyp(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Etikett:
		if didChange, err = k.addEtikett(kinder, mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Kasten:
		if didChange, err = k.addKasten(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Konfig:
		if didChange, err = k.setTransacted(kinder, ak.GetKonfigV0()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// switch kinder.Kennung.GetGattung() {
	// case gattung.Typ, gattung.Etikett, gattung.Kasten:
	// 	didChange = didChange && mutter != nil
	// }

	if didChange && (mutter != nil || mode.Contains(objekte_mode.ModeSchwanz)) {
		k.SetHasChanges(fmt.Sprintf("added: %s", kinder))
	}

	return
}

func (k *compiled) addTyp(
	b1 *sku.Transacted,
) (didChange bool, err error) {
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

	if didChange, err = iter.AddOrReplaceIfGreater(
		k.Typen,
		b,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
