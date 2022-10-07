package store_verzeichnisse

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

const digitWidth = 2
const pageCount = 1 << (digitWidth * 4)

type Zettelen struct {
	konfig.Konfig
	path string
	pool ZettelPool
	ioFactory
	pages [pageCount]*zettelenPageWithState
}

func MakeZettelen(
	k konfig.Konfig,
	s standort.Standort,
	f ioFactory,
) (i *Zettelen, err error) {
	i = &Zettelen{
		Konfig:    k,
		path:      s.DirVerzeichnisseZettelenNeue(),
		ioFactory: f,
		pool:      MakeZettelPool(),
	}

	for n, _ := range i.pages {
		i.pages[n] = makeZettelenPage(
			i,
			i.PathForPage(n),
			i.GetPageIndexKeyValue,
		)
	}

	return
}

func (i Zettelen) PathForPage(n int) (p string) {
	p = filepath.Join(i.path, fmt.Sprintf("%x", n))
	return
}

func (i Zettelen) ValidatePageIndex(n int) (err error) {
	switch {
	case n > pageCount:
		fallthrough

	case n < 0:
		err = errors.Errorf("expected page between 0 and %d, but got %d", pageCount-1, n)
		return
	}

	return
}

func (i Zettelen) PageForHinweis(h hinweis.Hinweis) (n int) {
	s := h.Sha()
	ss := s.String()[:digitWidth]

	var err error
	var n1 int64

	if n1, err = strconv.ParseInt(ss, 16, 64); err != nil {
		n1 = -1
	}

	n = int(n1)

	return
}

func (i Zettelen) PagesForZettelTransacted(zt zettel_transacted.Zettel) (ns []int) {
	ns = append(
		ns,
		i.PageForHinweis(zt.Named.Hinweis),
	)

	return
}

func (i *Zettelen) Flush() (err error) {
	errors.Print("flushing")

	for _, p := range i.pages {
		if err = p.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *Zettelen) Add(tz zettel_transacted.Zettel) (err error) {
	ns := i.PagesForZettelTransacted(tz)

	if len(ns) < 1 {
		err = errors.Errorf("expected at least one page value, but got none")
		return
	}

	n := ns[0]

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	z := i.pool.Get()
	z.PageSelection.Reason = PageSelectionReasonHinweis
	z.Transacted = tz
	z.EtikettenExpandedSorted = tz.Named.Stored.Zettel.Etiketten.Expanded().SortedString()
	z.EtikettenSorted = tz.Named.Stored.Zettel.Etiketten.SortedString()

	if err = p.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) GetPageIndexKeyValue(
	zt zettel_transacted.Zettel,
) (key string, value string) {
	key = zt.Named.Hinweis.String()
	value = fmt.Sprintf("%s.%s", zt.Schwanz, zt.Named.Stored.Sha)
	return
}

func (i *Zettelen) ReadHinweisSchwanzen(h hinweis.Hinweis) (tz zettel_transacted.Zettel, err error) {
	n := i.PageForHinweis(h)

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	if tz, err = p.ReadHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) IsSchwanz(z zettel_transacted.Zettel) (ok bool, err error) {
	n := i.PageForHinweis(z.Named.Hinweis)
	pi := i.pages[n]
	ok, err = pi.IsSchwanz(z)

	return
}

func (i *Zettelen) ReadMany(
	ws ...Writer,
) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(i.pages))

	w := writer{
		writers:    ws,
		ZettelPool: &i.pool,
	}

	for n, p := range i.pages {
		go func(n int, p *zettelenPageWithState) {
			defer wg.Done()

			if err = p.WriteZettelenTo(w); err != nil {
				err = errors.Wrap(err)
				return
			}
		}(n, p)
	}

	wg.Wait()

	return
}
