package store_verzeichnisse

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_verzeichnisse"
)

const DigitWidth = 2
const PageCount = 1 << (DigitWidth * 4)

type Zettelen struct {
	konfig.Konfig
	path string
	pool *zettel_verzeichnisse.Pool
	ioFactory
	pages [PageCount]*zettelenPageWithState
}

type pageId struct {
	index int
	path  string
}

func MakeZettelen(
	k konfig.Konfig,
	dir string,
	f ioFactory,
	p *zettel_verzeichnisse.Pool,
) (i *Zettelen, err error) {
	i = &Zettelen{
		Konfig:    k,
		path:      dir,
		ioFactory: f,
		pool:      p,
	}

	for n, _ := range i.pages {
		i.pages[n] = makeZettelenPage(
			f,
			i.PageIdForIndex(n),
			p,
		)
	}

	return
}

func (i Zettelen) PageIdForIndex(n int) (pid pageId) {
	pid.index = n
	pid.path = filepath.Join(i.path, fmt.Sprintf("%x", n))
	return
}

func (i Zettelen) ValidatePageIndex(n int) (err error) {
	switch {
	case n > PageCount:
		fallthrough

	case n < 0:
		err = errors.Errorf("expected page between 0 and %d, but got %d", PageCount-1, n)
		return
	}

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

func (i *Zettelen) Add(tz zettel_transacted.Zettel, v string) (err error) {
	var n int

	if n, err = i.PageForString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	z := i.pool.MakeZettel(tz)

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
	var n int

	if n, err = i.PageForHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (i *Zettelen) ReadMany(
	ws ...zettel_verzeichnisse.Writer,
) (err error) {
	wg := &sync.WaitGroup{}

	w := zettel_verzeichnisse.MakeWriterMulti(i.pool, ws...)

	for n, p := range i.pages {
		if n > 233 {
			continue
		}

		wg.Add(1)

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
