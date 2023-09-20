package store_verzeichnisse

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/pool"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/kilo/konfig"
)

const (
	DigitWidth = 2
	PageCount  = 1 << (DigitWidth * 4)
)

type Zettelen struct {
	erworben konfig.Compiled
	path     string
	pool     schnittstellen.Pool[transacted.Zettel, *transacted.Zettel]
	schnittstellen.VerzeichnisseFactory
	pages [PageCount]*Page
}

type pageId struct {
	index int
	path  string
}

func MakeZettelen(
	k konfig.Compiled,
	dir string,
	f schnittstellen.VerzeichnisseFactory,
	p schnittstellen.Pool[transacted.Zettel, *transacted.Zettel],
	fff PageDelegateGetter,
) (i *Zettelen, err error) {
	i = &Zettelen{
		erworben:             k,
		path:                 dir,
		VerzeichnisseFactory: f,
		pool:                 p,
	}

	for n := range i.pages {
		i.pages[n] = makeZettelenPage(
			f,
			i.PageIdForIndex(n),
			p,
			fff,
			// k.UseBestandsaufnahmeForVerzeichnisse,
			true,
		)
	}

	return
}

func (i Zettelen) Pool() schnittstellen.Pool[transacted.Zettel, *transacted.Zettel] {
	errors.TodoP4("rename to GetPool")
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
		err = errors.Errorf(
			"expected page between 0 and %d, but got %d",
			PageCount-1,
			n,
		)
		return
	}

	p = i.pages[n]

	return
}

func (i *Zettelen) SetNeedsFlush() {
	for _, p := range i.pages {
		p.State = StateChanged
	}
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

func (i *Zettelen) AddVerzeichnisse(
	tz sku.SkuLikePtr,
	v string,
) (err error) {
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

	if err = z.SetFromSkuLike(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = p.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) GetPageIndexKeyValue(
	zt transacted.Zettel,
) (key string, value string) {
	key = zt.Kennung.String()
	value = fmt.Sprintf("%s.%s", zt.GetTai(), zt.ObjekteSha)
	return
}

func (i *Zettelen) ReadMany(
	ws ...schnittstellen.FuncIter[*transacted.Zettel],
) (err error) {
	errors.TodoP3("switch to single writer and force callers to make chains")

	wg := &sync.WaitGroup{}
	ch := make(chan struct{}, PageCount)
	me := errors.MakeMulti()
	chDone := make(chan struct{})

	isDone := func() bool {
		select {
		case <-chDone:
			return true

		default:
			return false
		}
	}

	w := pool.MakePooledChain[transacted.Zettel](
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

				if err1 = p.Copy(w); err1 != nil {
					if isDone() {
						break
					}

					switch {
					case errors.IsTooManyOpenFiles(err1):
						<-openFileCh
						continue

					case iter.IsStopIteration(err1):
						break

					default:
						me.Add(err1)
						break
					}
				}

				break
			}
		}(n, p, ch)
	}

	wg.Wait()

	if me.Len() > 0 {
		err = me
	}

	return
}
