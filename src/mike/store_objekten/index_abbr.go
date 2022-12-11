package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type indexAbbrEncodableTridexes struct {
	Shas             *tridex.Tridex
	HinweisKopfen    *tridex.Tridex
	HinweisSchwanzen *tridex.Tridex
	Etiketten        *tridex.Tridex
}

type indexAbbr struct {
	ioFactory

	path string

	indexAbbrEncodableTridexes

	didRead    bool
	hasChanges bool
}

func newIndexAbbr(
	ioFactory ioFactory,
	p string,
) (i *indexAbbr, err error) {
	i = &indexAbbr{
		path:      p,
		ioFactory: ioFactory,
		indexAbbrEncodableTridexes: indexAbbrEncodableTridexes{
			Shas:             tridex.Make(),
			HinweisKopfen:    tridex.Make(),
			HinweisSchwanzen: tridex.Make(),
			Etiketten:        tridex.Make(),
		},
	}

	return
}

func (i *indexAbbr) Flush() (err error) {
	if !i.hasChanges {
		errors.Log().Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w1.Close)

	w := bufio.NewWriter(w1)

	defer errors.Deferred(&err, w.Flush)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.indexAbbrEncodableTridexes); err != nil {
		err = errors.Wrapf(err, "failed to write encoded kennung")
		return
	}

	return
}

func (i *indexAbbr) readIfNecessary() (err error) {
	errors.Log().Caller(1, "")

	if i.didRead {
		errors.Log().Print("already read")
		return
	}

	errors.Log().Print("reading")

	i.didRead = true

	var r1 io.ReadCloser

	if r1, err = i.ReadCloserVerzeichnisse(i.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, r1.Close)

	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	errors.Log().Print("starting decode")

	if err = dec.Decode(&i.indexAbbrEncodableTridexes); err != nil {
		errors.Log().Print("finished decode unsuccessfully")
		err = errors.Wrap(err)
		return
	}

	errors.Log().Print("finished decode successfully")

	return
}

func (i *indexAbbr) addStored(o gattung.Stored) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	i.indexAbbrEncodableTridexes.Shas.Add(o.AkteSha().String())
	i.indexAbbrEncodableTridexes.Shas.Add(o.ObjekteSha().String())

	return
}

func (i *indexAbbr) addZettelTransacted(zt zettel.Transacted) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	if err = i.addStored(&zt); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.indexAbbrEncodableTridexes.HinweisKopfen.Add(zt.Kennung().Kopf())
	i.indexAbbrEncodableTridexes.HinweisSchwanzen.Add(zt.Kennung().Schwanz())

	for _, e := range kennung.Expanded(zt.Objekte.Etiketten, kennung.ExpanderRight).Elements() {
		i.indexAbbrEncodableTridexes.Etiketten.Add(e.String())
	}

	return
}

func (i *indexAbbr) AbbreviateSha(s sha.Sha) (abbr string, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	abbr = i.indexAbbrEncodableTridexes.Shas.Abbreviate(s.String())

	return
}

func (i *indexAbbr) ExpandShaString(st string) (s sha.Sha, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	expanded := i.indexAbbrEncodableTridexes.Shas.Expand(st)

	if err = s.Set(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexAbbr) AbbreviateHinweis(h hinweis.Hinweis) (ha hinweis.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var kopf, schwanz string

	kopf = i.indexAbbrEncodableTridexes.HinweisKopfen.Abbreviate(h.Kopf())
	schwanz = i.indexAbbrEncodableTridexes.HinweisSchwanzen.Abbreviate(h.Schwanz())

	if kopf == "" || schwanz == "" {
		err = errors.Errorf("abbreviated kopf would be empty for %s", h)
		errors.Log().PrintDebug(i.indexAbbrEncodableTridexes.HinweisKopfen)
		return
	}

	if schwanz == "" {
		err = errors.Errorf("abbreviated schwanz would be empty for %s", h)
		return
	}

	if ha, err = hinweis.MakeKopfUndSchwanz(kopf, schwanz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexAbbr) ExpandHinweisString(s string) (h hinweis.Hinweis, err error) {
	errors.Log().Print(s)

	var ha hinweis.Hinweis

	if ha, err = hinweis.Make(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.ExpandHinweis(ha)
}

func (i *indexAbbr) ExpandHinweis(hAbbr hinweis.Hinweis) (h hinweis.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kopf := i.indexAbbrEncodableTridexes.HinweisKopfen.Expand(hAbbr.Kopf())
	schwanz := i.indexAbbrEncodableTridexes.HinweisSchwanzen.Expand(hAbbr.Schwanz())

	if h, err = hinweis.MakeKopfUndSchwanz(kopf, schwanz); err != nil {
		err = errors.Wrapf(err, "{Abbreviated: '%s'}", hAbbr)
		return
	}

	return
}

func (i *indexAbbr) ExpandEtikettString(s string) (e kennung.Etikett, err error) {
	errors.Log().Print(s)

	if e, err = kennung.MakeEtikett(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.ExpandEtikett(e)
}

func (i *indexAbbr) ExpandEtikett(eAbbr kennung.Etikett) (e kennung.Etikett, err error) {
	errors.Log().Print(eAbbr)

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ex := i.indexAbbrEncodableTridexes.Etiketten.Expand(eAbbr.String())

	if ex == "" {
		//TODO should try to use the expansion if possible
		ex = eAbbr.String()
	}

	if err = e.Set(ex); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
