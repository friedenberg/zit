package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

type indexZettelenTails struct {
	umwelt *umwelt.Umwelt
	path   string
	ioFactory
	zettelen   zettel_transacted.Set
	didRead    bool
	hasChanges bool
}

func newIndexZettelenTails(
	u *umwelt.Umwelt,
	p string,
	f ioFactory,
) (i *indexZettelenTails, err error) {
	i = &indexZettelenTails{
		umwelt:    u,
		path:      p,
		ioFactory: f,
		zettelen:  zettel_transacted.MakeSetHinweisTransacted(),
	}

	return
}

func (i *indexZettelenTails) Flush() (err error) {
	if !i.hasChanges {
		errors.Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Error(err)
		return
	}

	defer stdprinter.PanicIfError(w1.Close)

	w := bufio.NewWriter(w1)

	defer stdprinter.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	err = i.zettelen.Each(
		func(tz zettel_transacted.Zettel) (err error) {
			if err = enc.Encode(tz); err != nil {
				err = errors.Wrapped(err, "failed to write zettel: %s", tz.Named)
				return
			}

			return
		},
	)

	if err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (i *indexZettelenTails) readIfNecessary() (err error) {
	if i.didRead {
		return
	}

	i.didRead = true

	var r1 io.ReadCloser

	if r1, err = i.ReadCloserVerzeichnisse(i.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Error(err)
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
				err = errors.Error(err)
				return
			}
		}

		i.zettelen.Add(tz)
	}

	return
}

func (i *indexZettelenTails) add(tz zettel_transacted.Zettel) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	i.hasChanges = true

	i.zettelen.Add(tz)

	return
}

func (i *indexZettelenTails) Read(h hinweis.Hinweis) (tz zettel_transacted.Zettel, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if tz, ok = i.zettelen.Get(h); !ok {
		err = ErrNotFound{Id: h}
		return
	}

	return
}

func (i *indexZettelenTails) ZettelenSchwanzen(
	qs ...zettel_named.NamedFilter,
) (tzs map[hinweis.Hinweis]zettel_transacted.Zettel, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	tzs = make(map[hinweis.Hinweis]zettel_transacted.Zettel)

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

			tzs[tz.Named.Hinweis] = tz

			return
		},
	)

	if err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (i *indexZettelenTails) shouldIncludeTransacted(tz zettel_transacted.Zettel) bool {
	if i.umwelt.Konfig.IncludeHidden {
		return true
	}

	prefixes := tz.Named.Stored.Zettel.Etiketten.Expanded(etikett.ExpanderRight{})

	for tn, tv := range i.umwelt.Konfig.Tags {
		if !tv.Hide {
			continue
		}

		if prefixes.ContainsString(tn) {
			return false
		}
	}

	return true
}
