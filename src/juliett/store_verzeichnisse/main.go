package store_verzeichnisse

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/india/zettel_verzeichnisse"
)

const DigitWidth = 2
const PageCount = 1 << (DigitWidth * 4)

type Zettelen struct {
	konfig.Konfig
	path string
	pool *zettel_verzeichnisse.Pool
	ioFactory
	pages [PageCount]*Page
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

func (i Zettelen) Pool() *zettel_verzeichnisse.Pool {
	return i.pool
}

func (i Zettelen) PageIdForIndex(n int) (pid pageId) {
	pid.index = n
	pid.path = filepath.Join(i.path, fmt.Sprintf("%x", n))
	return
}

func (i Zettelen) GetPage(n int) (p *Page, err error) {
	switch {
	case n > PageCount:
		fallthrough

	case n < 0:
		err = errors.Errorf("expected page between 0 and %d, but got %d", PageCount-1, n)
		return
	}

	p = i.pages[n]

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

	var p *Page

	if p, err = i.GetPage(n); err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (i *Zettelen) ReadMany(
	ws ...zettel_verzeichnisse.Writer,
) (err error) {
	wg := &sync.WaitGroup{}
	ch := make(chan struct{}, PageCount)

	w1 := zettel_verzeichnisse.MakeWriterChain(ws...)
	w := zettel_verzeichnisse.MakeWriterChainIgnoreEOF(w1, i.pool)

	for n, p := range i.pages {
		wg.Add(1)

		go func(n int, p *Page, openFileCh chan struct{}) {
			defer wg.Done()
			defer func(c chan<- struct{}) {
				openFileCh <- struct{}{}
			}(openFileCh)

			for {
				if err = p.WriteZettelenTo(w); err != nil {
					if errors.IsTooManyOpenFiles(err) {
						<-openFileCh
						continue
					}

					//TODO hand back error
					err = errors.Wrap(err)
					errors.Err().Print(err)
					// return
				}

				break
			}

		}(n, p, ch)
	}

	wg.Wait()

	return
}
