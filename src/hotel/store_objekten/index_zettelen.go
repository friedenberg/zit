package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type indexZettelen struct {
	umwelt *umwelt.Umwelt
	path   string
	ioFactory
	zettelen      map[sha.Sha]stored_zettel.Transacted
	hinweisen     map[hinweis.Hinweis]stored_zettel.SetTransacted
	akten         map[sha.Sha]stored_zettel.SetTransacted
	bezeichnungen map[string]stored_zettel.SetTransacted
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
		zettelen:      make(map[sha.Sha]stored_zettel.Transacted),
		hinweisen:     make(map[hinweis.Hinweis]stored_zettel.SetTransacted),
		akten:         make(map[sha.Sha]stored_zettel.SetTransacted),
		bezeichnungen: make(map[string]stored_zettel.SetTransacted),
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

	for h, tz := range i.zettelen {
		if err = enc.Encode(tz); err != nil {
			err = errors.Wrapped(err, "failed to write zettel: [%s %s]", h, tz.Sha)
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
	i.zettelen[tz.Sha] = tz

	var set stored_zettel.SetTransacted
	var ok bool

	if set, ok = i.hinweisen[tz.Hinweis]; !ok {
		set = stored_zettel.MakeSetTransacted()
	}

	set.Add(tz)
	i.hinweisen[tz.Hinweis] = set

	if set, ok = i.akten[tz.Zettel.Akte]; !ok {
		set = stored_zettel.MakeSetTransacted()
	}

	set.Add(tz)
	i.akten[tz.Zettel.Akte] = set

	key := strings.ToLower(tz.Zettel.Bezeichnung.String())

	if set, ok = i.bezeichnungen[key]; !ok {
		set = stored_zettel.MakeSetTransacted()
	}

	set.Add(tz)

	i.bezeichnungen[key] = set
}

func (i *indexZettelen) Add(tz stored_zettel.Transacted) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	i.hasChanges = true

	i.addNoRead(tz)

	return
}

func (i *indexZettelen) ReadHinweis(h hinweis.Hinweis) (tzs stored_zettel.SetTransacted, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if tzs, ok = i.hinweisen[h]; !ok {
		err = ErrNotFound{Id: h}
		return
	}

	return
}

func (i *indexZettelen) ReadBezeichnung(s string) (tzs stored_zettel.SetTransacted, err error) {
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

func (i *indexZettelen) ReadAkteSha(s sha.Sha) (tzs stored_zettel.SetTransacted, err error) {
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

func (i *indexZettelen) ReadZettelSha(s sha.Sha) (tz stored_zettel.Transacted, err error) {
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
