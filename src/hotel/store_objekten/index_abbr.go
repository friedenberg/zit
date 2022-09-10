package store_objekten

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/tridex"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type indexAbbrEncodableTridexes struct {
	Shas             *tridex.Tridex
	HinweisKopfen    *tridex.Tridex
	HinweisSchwanzen *tridex.Tridex
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
		},
	}

	return
}

func (i *indexAbbr) Flush() (err error) {
	errCtx := errors.Ctx{}

	defer func() {
		err = errCtx.Error()
	}()

	if !i.hasChanges {
		errors.Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, errCtx.Err = i.WriteCloserVerzeichnisse(i.path); !errCtx.IsEmpty() {
		errCtx.Wrap()
		return
	}

	defer errCtx.Defer(w1.Close)

	w := bufio.NewWriter(w1)

	defer errCtx.Defer(w.Flush)

	enc := gob.NewEncoder(w)

	if errCtx.Err = enc.Encode(i.indexAbbrEncodableTridexes); !errCtx.IsEmpty() {
		errCtx.Wrapf("failed to write encoded kennung")
		return
	}

	return
}

func (i *indexAbbr) readIfNecessary() (err error) {
	errors.Caller(1, "")
	errCtx := errors.Ctx{}

	defer func() {
		err = errCtx.Error()
	}()

	if i.didRead {
		errors.Print("already read")
		return
	}

	errors.Print("reading")

	i.didRead = true

	var r1 io.ReadCloser

	if r1, errCtx.Err = i.ReadCloserVerzeichnisse(i.path); !errCtx.IsEmpty() {
		if errors.IsNotExist(errCtx.Err) {
			errCtx.ClearErr()
		} else {
			errCtx.Wrap()
		}

		return
	}

	defer errCtx.Defer(r1.Close)

	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	errors.Print("starting decode")

	if errCtx.Err = dec.Decode(&i.indexAbbrEncodableTridexes); !errCtx.IsEmpty() {
		errors.Print("finished decode unsuccessfully")
		errCtx.Wrap()
		return
	}

	errors.Print("finished decode successfully")

	return
}

func (i *indexAbbr) addZettelTransacted(zt zettel_transacted.Zettel) (err error) {
	ctx := errors.Ctx{}

	defer func() {
		err = ctx.Error()
	}()

	i.hasChanges = true

	if ctx.Err = i.readIfNecessary(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	i.indexAbbrEncodableTridexes.Shas.Add(zt.Named.Stored.Sha.String())
	i.indexAbbrEncodableTridexes.Shas.Add(zt.Named.Stored.Zettel.Akte.String())
	i.indexAbbrEncodableTridexes.HinweisKopfen.Add(zt.Named.Hinweis.Kopf())
	i.indexAbbrEncodableTridexes.HinweisSchwanzen.Add(zt.Named.Hinweis.Schwanz())

	return
}

func (i *indexAbbr) AbbreviateSha(s sha.Sha) (abbr string, err error) {
	ctx := errors.Ctx{}

	defer func() {
		err = ctx.Error()
	}()

	if ctx.Err = i.readIfNecessary(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	abbr = i.indexAbbrEncodableTridexes.Shas.Abbreviate(s.String())

	return
}

func (i *indexAbbr) ExpandShaString(st string) (s sha.Sha, err error) {
	ctx := errors.Ctx{}

	defer func() {
		err = ctx.Error()
	}()

	if ctx.Err = i.readIfNecessary(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	expanded := i.indexAbbrEncodableTridexes.Shas.Expand(st)

	if ctx.Err = s.Set(expanded); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	return
}

func (i *indexAbbr) AbbreviateHinweis(h hinweis.Hinweis) (abbr string, err error) {
	ctx := errors.Ctx{}

	defer func() {
		err = ctx.Error()
	}()

	if ctx.Err = i.readIfNecessary(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	var kopf, schwanz string

	kopf = i.indexAbbrEncodableTridexes.HinweisKopfen.Abbreviate(h.Kopf())
	schwanz = i.indexAbbrEncodableTridexes.HinweisSchwanzen.Abbreviate(h.Schwanz())
	// kopf = h.Kopf()
	// schwanz = h.Schwanz()
	abbr = fmt.Sprintf("%s/%s", kopf, schwanz)

	return
}

func (i *indexAbbr) ExpandHinweisString(st string) (h hinweis.Hinweis, err error) {
	ctx := errors.Ctx{}

	defer func() {
		err = ctx.Error()
	}()

	if ctx.Err = i.readIfNecessary(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	// expanded := i.indexAbbrEncodableTridexes.Hinweis.Expand(st)

	// if ctx.Err = s.Set(expanded); !ctx.IsEmpty() {
	// 	ctx.Wrap()
	// 	return
	// }

	return
}
