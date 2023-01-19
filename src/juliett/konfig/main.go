package konfig

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"

	pkg_angeboren "github.com/friedenberg/zit/src/alfa/angeboren"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/etikett"
	"github.com/friedenberg/zit/src/india/typ"
)

var (
	typExpander kennung.Expander
)

func init() {
	typExpander = kennung.MakeExpanderRight(`-`)
}

type angeboren = pkg_angeboren.Konfig

type Compiled struct {
	cli
	compiled
	angeboren
}

// TODO-P4 rename this
func (a *Compiled) ResetWithInner(b compiled) {
	a.compiled = b
}

type cli = erworben.Cli

type compiled struct {
	hasChanges bool

	Sku sku.Transacted[kennung.Konfig, *kennung.Konfig]

	erworben.Akte

	//Etiketten
	EtikettenHidden     []string
	EtikettenToAddToNew []string
	Etiketten           etikettSet

	//Typen
	ExtensionsToTypen map[string]string
	DefaultTyp        typ.Transacted
	Typen             typSet
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

	if err = c.loadKonfigAngeboren(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.loadKonfigErworben(s); err != nil {
		err = errors.Wrap(err)
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

	if err = kc.Etiketten.Each(
		func(ct *etikett.Transacted) (err error) {
			tn := ct.Sku.Kennung.String()
			tv := ct.Objekte.Akte

			switch {
			case tv.Hide:
				kc.EtikettenHidden = append(kc.EtikettenHidden, tn)

			case tv.AddToNewZettels:
				kc.EtikettenToAddToNew = append(kc.EtikettenToAddToNew, tn)
			}

			kc.applyExpandedEtikett(ct)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(kc.EtikettenHidden, func(i, j int) bool {
		return kc.EtikettenHidden[i] < kc.EtikettenHidden[j]
	})

	sort.Slice(kc.EtikettenToAddToNew, func(i, j int) bool {
		return kc.EtikettenToAddToNew[i] < kc.EtikettenToAddToNew[j]
	})

	if err = kc.Typen.Each(
		func(ct *typ.Transacted) (err error) {
			fe := ct.Objekte.Akte.FileExtension

			if fe != "" {
				kc.ExtensionsToTypen[fe] = ct.Sku.Kennung.String()
			}

			if ct == nil {
				errors.Todo("determine why any Typen might be nil")
				return
			}

			kc.applyExpandedTyp(*ct)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *Compiled) Flush(s standort.Standort) (err error) {
	if !kc.hasChanges {
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

	defer errors.Deferred(&err, f.Close)

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
func (c compiled) GetSortedTypenExpanded(v string) (expandedActual []*typ.Transacted) {
	expandedMaybe := collections.MakeMutableValueSet[collections.StringValue, *collections.StringValue]()
	sa := collections.MakeFuncSetString[collections.StringValue, *collections.StringValue](expandedMaybe)
	typExpander.Expand(sa, v)
	expandedActual = make([]*typ.Transacted, 0)

	expandedMaybe.Each(
		func(v collections.StringValue) (err error) {
			ct, ok := c.Typen.Get(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, ct)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return expandedActual[i].Sku.Kennung.Len() > expandedActual[j].Sku.Kennung.Len()
	})

	return
}

func (c compiled) GetSortedEtikettenExpanded(
	v string,
) (expandedActual []*etikett.Transacted) {
	expandedMaybe := collections.MakeMutableValueSet[collections.StringValue, *collections.StringValue]()
	sa := collections.MakeFuncSetString[collections.StringValue, *collections.StringValue](expandedMaybe)
	typExpander.Expand(sa, v)
	expandedActual = make([]*etikett.Transacted, 0)

	expandedMaybe.Each(
		func(v collections.StringValue) (err error) {
			ct, ok := c.Etiketten.Get(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, ct)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return expandedActual[i].Sku.Kennung.Len() > expandedActual[j].Sku.Kennung.Len()
	})

	return
}

func (kc compiled) IsInlineTyp(k kennung.Typ) (isInline bool) {
	tc := kc.GetTyp(k)

	if tc == nil {
		return
	}

	isInline = tc.Objekte.Akte.InlineAkte

	return
}

func (kc compiled) GetTyp(k kennung.Typ) (ct *typ.Transacted) {
	expandedActual := kc.GetSortedTypenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
	}

	return
}

func (kc compiled) GetEtikett(k kennung.Etikett) (ct *etikett.Transacted) {
	expandedActual := kc.GetSortedEtikettenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
	}

	return
}

func (k *compiled) SetTransacted(
	kt *erworben.Transacted,
) {
	k.hasChanges = true
	k.Sku = kt.Sku
	k.Akte = kt.Objekte.Akte
}

func (k *compiled) AddTyp(
	ct *typ.Transacted,
) {
	if ct.Objekte.Akte.Actions == nil {
		errors.TodoP0("actions were nil: %s", ct.Sku)
		return
	}

	if ct.Objekte.Akte.EtikettenRules == nil {
		errors.TodoP0("etiketten rules were nil: %s", ct.Sku)
		return
	}

	k.hasChanges = true
	// collections.AddIfGreater(k.Typen, ct)
	m := k.Typen.Elements()
	m = append(m, ct)
	k.Typen = makeCompiledTypSetFromSlice(m)

	return
}

func (k *compiled) AddEtikett(
	ct *etikett.Transacted,
) {
	k.hasChanges = true
	m := k.Etiketten.Elements()
	m = append(m, ct)
	k.Etiketten = makeCompiledEtikettSetFromSlice(m)

	return
}

func (c *compiled) applyExpandedTyp(ct typ.Transacted) {
	expandedActual := c.GetSortedTypenExpanded(ct.Sku.Kennung.String())

	for _, ex := range expandedActual {
		ct.Objekte.Akte.Merge(&ex.Objekte.Akte)
	}
}

func (c *compiled) applyExpandedEtikett(ct *etikett.Transacted) {
	expandedActual := c.GetSortedEtikettenExpanded(ct.Sku.Kennung.String())

	for _, ex := range expandedActual {
		ct.Objekte.Akte.Merge(&ex.Objekte.Akte)
	}
}
