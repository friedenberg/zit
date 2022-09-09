package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/trie"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type indexAbbrEncodableTries struct {
	Shas             *trie.Trie
	HinweisKopfen    *trie.Trie
	HinweisSchwanzen *trie.Trie
}

type indexAbbr struct {
	*umwelt.Umwelt
	ioFactory

	path string

	indexAbbrEncodableTries

	didRead    bool
	hasChanges bool
}

func newIndexAbbr(
	u *umwelt.Umwelt,
	ioFactory ioFactory,
	p string,
) (i *indexAbbr, err error) {
	i = &indexAbbr{
		Umwelt:    u,
		path:      p,
		ioFactory: ioFactory,
		indexAbbrEncodableTries: indexAbbrEncodableTries{
			Shas:             trie.Make(),
			HinweisKopfen:    trie.Make(),
			HinweisSchwanzen: trie.Make(),
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

	if errCtx.Err = enc.Encode(i.indexAbbrEncodableTries); !errCtx.IsEmpty() {
		errCtx.Wrapf("failed to write encoded kennung")
		return
	}

	return
}

func (i *indexAbbr) readIfNecessary() (err error) {
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

	if errCtx.Err = dec.Decode(&i.indexAbbrEncodableTries); !errCtx.IsEmpty() {
		errCtx.Wrap()
		return
	}

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

	i.indexAbbrEncodableTries.Shas.Add(zt.Named.Stored.Sha.String())
	i.indexAbbrEncodableTries.Shas.Add(zt.Named.Stored.Zettel.Akte.String())
	i.indexAbbrEncodableTries.HinweisKopfen.Add(zt.Named.Hinweis.Kopf())
	i.indexAbbrEncodableTries.HinweisSchwanzen.Add(zt.Named.Hinweis.Schwanz())

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

	abbr = i.indexAbbrEncodableTries.Shas.Abbreviate(s.String())

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

	expanded := i.indexAbbrEncodableTries.Shas.Expand(st)

	if ctx.Err = s.Set(expanded); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	return
}
