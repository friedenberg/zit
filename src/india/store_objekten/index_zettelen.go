package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/bravo/typ"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

type indexZettelen struct {
	umwelt *umwelt.Umwelt
	path   string
	ioFactory
	zettelen      map[sha.Sha]zettel_transacted.Set
	hinweisen     map[hinweis.Hinweis]zettel_transacted.Set
	akten         map[sha.Sha]zettel_transacted.Set
	bezeichnungen map[string]zettel_transacted.Set
	typen         map[typ.Typ]zettel_transacted.Set
	didRead       bool
	hasChanges    bool
}

func newIndexZettelen(
	u *umwelt.Umwelt,
	p string,
	f ioFactory,
) (i *indexZettelen, err error) {
	i = &indexZettelen{
		umwelt:        u,
		path:          p,
		ioFactory:     f,
		zettelen:      make(map[sha.Sha]zettel_transacted.Set),
		hinweisen:     make(map[hinweis.Hinweis]zettel_transacted.Set),
		akten:         make(map[sha.Sha]zettel_transacted.Set),
		bezeichnungen: make(map[string]zettel_transacted.Set),
		typen:         make(map[typ.Typ]zettel_transacted.Set),
	}

	return
}

func (i *indexZettelen) Flush() (err error) {
	if !i.hasChanges {
		logz.Print("no changes")
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

	for _, st := range i.zettelen {
		err = st.Each(
			func(tz zettel_transacted.Zettel) (err error) {
				if err = enc.Encode(tz); err != nil {
					err = errors.Wrapped(err, "failed to write zettel: [%s %s]", tz.Named.Hinweis, tz.Named.Stored.Sha)
					return
				}

				return
			},
		)

		if err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (i *indexZettelen) readIfNecessary() (err error) {
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

		i.addNoRead(tz)
	}

	return
}

func (i *indexZettelen) addNoRead(tz zettel_transacted.Zettel) {
	{
		var set zettel_transacted.Set
		var ok bool

		if set, ok = i.zettelen[tz.Named.Stored.Sha]; !ok {
			set = zettel_transacted.MakeSetUnique(1)
		}

		set.Add(tz)
		i.zettelen[tz.Named.Stored.Sha] = set
	}

	{
		var set zettel_transacted.Set
		var ok bool

		if set, ok = i.hinweisen[tz.Named.Hinweis]; !ok {
			set = zettel_transacted.MakeSetUnique(1)
		}

		set.Add(tz)
		i.hinweisen[tz.Named.Hinweis] = set
	}

	akteSha := tz.Named.Stored.Zettel.Akte

	if !akteSha.IsNull() {
		var set zettel_transacted.Set
		var ok bool

		if set, ok = i.akten[tz.Named.Stored.Zettel.Akte]; !ok {
			set = zettel_transacted.MakeSetUnique(1)
		}

		set.Add(tz)
		i.akten[tz.Named.Stored.Zettel.Akte] = set
	}

	bezKey := strings.ToLower(tz.Named.Stored.Zettel.Bezeichnung.String())
	if bezKey != "" {

		var set zettel_transacted.Set
		var ok bool

		if set, ok = i.bezeichnungen[bezKey]; !ok {
			set = zettel_transacted.MakeSetUnique(1)
		}

		set.Add(tz)
		i.bezeichnungen[bezKey] = set
	}

	{
		var set zettel_transacted.Set
		var ok bool

		if set, ok = i.typen[tz.Named.Stored.Zettel.Typ]; !ok {
			set = zettel_transacted.MakeSetUnique(1)
		}

		set.Add(tz)
		i.typen[tz.Named.Stored.Zettel.Typ] = set
	}

	return
}

func (i *indexZettelen) add(tz zettel_transacted.Zettel) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	i.hasChanges = true

	i.addNoRead(tz)

	return
}

func (i *indexZettelen) ReadHinweis(h hinweis.Hinweis) (mst zettel_transacted.Set, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if mst, ok = i.hinweisen[h]; !ok {
		err = ErrNotFound{Id: h}
		return
	}

	return
}

func (i *indexZettelen) ReadBezeichnung(s string) (tzs zettel_transacted.Set, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if tzs, ok = i.bezeichnungen[s]; !ok {
		err = ErrNotFound{Id: stringId(s)}
		return
	}

	return
}

func (i *indexZettelen) ReadAkteSha(s sha.Sha) (tzs zettel_transacted.Set, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if tzs, ok = i.akten[s]; !ok {
		err = ErrNotFound{Id: s}
		return
	}

	return
}

func (i *indexZettelen) ReadZettelSha(s sha.Sha) (tz zettel_transacted.Set, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if tz, ok = i.zettelen[s]; !ok {
		err = ErrNotFound{Id: s}
		return
	}

	return
}

func (i *indexZettelen) ReadTyp(t typ.Typ) (tzs zettel_transacted.Set, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if tzs, ok = i.typen[t]; !ok {
		err = ErrNotFound{Id: t}
		return
	}

	return
}
