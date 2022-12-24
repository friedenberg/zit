package store_verzeichnisse

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

const DigitWidth = 2
const PageCount = 1 << (DigitWidth * 4)

type Zettelen struct {
	konfig konfig_compiled.Compiled
	path   string
	pool   *zettel.PoolVerzeichnisse
	ioFactory
	pages [PageCount]*Page
}

type pageId struct {
	index int
	path  string
}

func MakeZettelen(
	k konfig_compiled.Compiled,
	dir string,
	f ioFactory,
	p *zettel.PoolVerzeichnisse,
	fff ZettelVerzeichnisseWriterGetter,
) (i *Zettelen, err error) {
	i = &Zettelen{
		konfig:    k,
		path:      dir,
		ioFactory: f,
		pool:      p,
	}

	for n := range i.pages {
		i.pages[n] = makeZettelenPage(
			f,
			i.PageIdForIndex(n),
			p,
			fff,
		)
	}

	return
}

func (i Zettelen) Pool() *zettel.PoolVerzeichnisse {
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
	errors.Log().Print("flushing")

	for _, p := range i.pages {
		if err = p.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

// TODO-P2 switch to pointer
func (i *Zettelen) Add(tz *zettel.Transacted, v string) (err error) {
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

	z := i.pool.Get()
	z.Reset(tz)

	if err = p.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) GetPageIndexKeyValue(
	zt zettel.Transacted,
) (key string, value string) {
	key = zt.Kennung().String()
	value = fmt.Sprintf("%s.%s", zt.Sku.Schwanz, zt.Sku.Sha)
	return
}

func (i *Zettelen) ReadMany(
	//TODO switch to single writer and force callers to make chains
	ws ...collections.WriterFunc[*zettel.Transacted],
) (err error) {
	wg := &sync.WaitGroup{}
	ch := make(chan struct{}, PageCount)
	chErr := make(chan error)
	chDone := make(chan struct{})

	isDone := func() bool {
		select {
		case <-chDone:
			return true

		default:
			return false
		}
	}

	w := collections.MakePooledChain[zettel.Transacted](
		i.pool,
		ws...,
	)

	for n, p := range i.pages {
		wg.Add(1)

		go func(n int, p *Page, openFileCh chan struct{}) {
			defer wg.Done()
			defer func(c chan<- struct{}) {
				openFileCh <- struct{}{}
			}(openFileCh)

			for {
				if isDone() {
					break
				}

				var err1 error

				if err1 = p.WriteZettelenTo(w); err1 != nil {
					if isDone() {
						break
					}

					switch {
					case errors.IsTooManyOpenFiles(err1):
						<-openFileCh
						continue

					case errors.IsEOF(err1):
						break

					default:
						chErr <- err1
						break
					}
				}

				break
			}

		}(n, p, ch)
	}

	go func() {
		err = <-chErr
		close(chDone)
	}()

	wg.Wait()

	return
}
