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
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/hotel/collections"
)

type indexZettelen struct {
	umwelt *umwelt.Umwelt
	path   string
	ioFactory
	zettelen      map[sha.Sha]collections.SetTransacted
	hinweisen     map[hinweis.Hinweis]collections.SetTransacted
	akten         map[sha.Sha]collections.SetTransacted
	bezeichnungen map[string]collections.SetTransacted
	typen         map[typ.Typ]collections.SetTransacted
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
		zettelen:      make(map[sha.Sha]collections.SetTransacted),
		hinweisen:     make(map[hinweis.Hinweis]collections.SetTransacted),
		akten:         make(map[sha.Sha]collections.SetTransacted),
		bezeichnungen: make(map[string]collections.SetTransacted),
		typen:         make(map[typ.Typ]collections.SetTransacted),
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
			func(tz stored_zettel.Transacted) (err error) {
				if err = enc.Encode(tz); err != nil {
					err = errors.Wrapped(err, "failed to write zettel: [%s %s]", tz.Hinweis, tz.Stored.Sha)
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
		var tz stored_zettel.Transacted

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

func (i *indexZettelen) addNoRead(tz stored_zettel.Transacted) {
	{
		var set collections.SetTransacted
		var ok bool

		if set, ok = i.zettelen[tz.Stored.Sha]; !ok {
			set = collections.MakeSetUniqueTransacted()
		}

		set.Add(tz)
		i.zettelen[tz.Stored.Sha] = set
	}

	{
		var set collections.SetTransacted
		var ok bool

		if set, ok = i.hinweisen[tz.Hinweis]; !ok {
			set = collections.MakeSetUniqueTransacted()
		}

		set.Add(tz)
		i.hinweisen[tz.Hinweis] = set
	}

	{
		var set collections.SetTransacted
		var ok bool

		if set, ok = i.akten[tz.Zettel.Akte]; !ok {
			set = collections.MakeSetUniqueTransacted()
		}

		set.Add(tz)
		i.akten[tz.Zettel.Akte] = set
	}

	{
		key := strings.ToLower(tz.Zettel.Bezeichnung.String())
		var set collections.SetTransacted
		var ok bool

		if set, ok = i.bezeichnungen[key]; !ok {
			set = collections.MakeSetUniqueTransacted()
		}

		set.Add(tz)
		i.bezeichnungen[key] = set
	}

	{
		var set collections.SetTransacted
		var ok bool

		if set, ok = i.typen[tz.Zettel.Typ]; !ok {
			set = collections.MakeSetUniqueTransacted()
		}

		set.Add(tz)
		i.typen[tz.Zettel.Typ] = set
	}

	return
}

func (i *indexZettelen) add(tz stored_zettel.Transacted) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	i.hasChanges = true

	i.addNoRead(tz)

	return
}

func (i *indexZettelen) ReadHinweis(h hinweis.Hinweis) (mst collections.SetTransacted, err error) {
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

func (i *indexZettelen) ReadBezeichnung(s string) (tzs collections.SetTransacted, err error) {
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

func (i *indexZettelen) ReadAkteSha(s sha.Sha) (tzs collections.SetTransacted, err error) {
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

func (i *indexZettelen) ReadZettelSha(s sha.Sha) (tz collections.SetTransacted, err error) {
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

func (i *indexZettelen) ReadTyp(t typ.Typ) (tzs collections.SetTransacted, err error) {
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
