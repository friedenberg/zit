package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/verzeichnisse"
)

type indexZettelenTails struct {
	konfig.Konfig
	path string
	ioFactory
	zettelen   zettel_transacted.Set
	didRead    bool
	hasChanges bool

	*verzeichnisse.Zettelen
}

func newIndexZettelenTails(
	k konfig.Konfig,
	p string,
	f ioFactory,
	pool zettel_transacted.Pool,
) (i *indexZettelenTails, err error) {
	i = &indexZettelenTails{
		Konfig:    k,
		path:      p,
		ioFactory: f,
		zettelen:  zettel_transacted.MakeSetHinweis(0),
	}

	var s standort.Standort

	if s, err = standort.Make(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	if i.Zettelen, err = verzeichnisse.MakeZettelen(k, s, f, pool); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexZettelenTails) Flush() (err error) {
	if err = i.Zettelen.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !i.hasChanges {
		errors.Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(w1.Close)

	w := bufio.NewWriter(w1)

	defer errors.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	err = i.zettelen.Each(
		func(tz zettel_transacted.Zettel) (err error) {
			if err = enc.Encode(tz); err != nil {
				err = errors.Wrapf(err, "failed to write zettel: %s", tz.Named)
				return
			}

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexZettelenTails) readIfNecessary() (err error) {
	if i.didRead {
		errors.Print("already read")
		return
	}

	errors.Print("read start")
	defer errors.Print("read end")

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

	defer r1.Close()

	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	for {
		var tz zettel_transacted.Zettel

		if err = dec.Decode(&tz); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		i.zettelen.Add(tz)
	}

	return
}

func (i *indexZettelenTails) add(tz zettel_transacted.Zettel) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	i.zettelen.Add(tz)

	if err = i.Zettelen.Add(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexZettelenTails) IsSchwanz(zt zettel_transacted.Zettel) (ok bool, err error) {
	var zta zettel_transacted.Zettel

	if zta, err = i.Read(zt.Named.Hinweis); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	ok = zta.Named.Equals(zt.Named)

	return
}

func (i *indexZettelenTails) Read(h hinweis.Hinweis) (tz zettel_transacted.Zettel, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ok := false

	if tz, ok = i.zettelen.Get(h); !ok {
		err = ErrNotFound{Id: h}
		return
	}

	tz.Named.Stored.Zettel.Etiketten = tz.Named.Stored.Zettel.Etiketten.Copy()
	var tz1 zettel_transacted.Zettel

	if tz1, err = i.Zettelen.ReadHinweisSchwanzen(h); err != nil {
		err = errors.Wrap(err)
		return
	} else if !tz1.Named.Equals(tz.Named) {
		err = errors.Errorf("ZettelenNeue had different zettel:\nneue: %s\nold: %s", tz1, tz)
		return
	}

	return
}

func (i *indexZettelenTails) ReadManySchwanzen(
	ws ...verzeichnisse.Writer,
) (err error) {
	return i.Zettelen.ReadMany(
		append(
			[]verzeichnisse.Writer{
				i.ZettelWriterSchwanzenOnly(),
				i.ZettelWriterFilterHidden(),
			},
			ws...,
		)...,
	)
}

func (i *indexZettelenTails) ZettelenSchwanzen(
	qs ...zettel_named.NamedFilter,
) (zts zettel_transacted.Set, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	zts = zettel_transacted.MakeSetUnique(0)

	err = i.zettelen.Each(
		func(tz zettel_transacted.Zettel) (err error) {
			for _, q := range qs {
				if !q.IncludeNamedZettel(tz.Named) {
					return
				}
			}

			if !i.shouldIncludeTransacted(tz) {
				return
			}

			zts.Add(tz)

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexZettelenTails) shouldIncludeTransacted(tz zettel_transacted.Zettel) bool {
	if i.IncludeHidden {
		return true
	}

	prefixes := tz.Named.Stored.Zettel.Etiketten.Expanded(etikett.ExpanderRight{})

	for tn, tv := range i.Tags {
		if !tv.Hide {
			continue
		}

		if prefixes.ContainsString(tn) {
			return false
		}
	}

	return true
}
