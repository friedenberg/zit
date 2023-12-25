package store_verzeichnisse

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/pool"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

type State int

const (
	StateUnread = State(iota)
	StateChanged
)

type PageDelegate interface {
	ShouldAddVerzeichnisse(*sku.Transacted) error
	ShouldFlushVerzeichnisse(*sku.Transacted) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(uint8) PageDelegate
}

const (
	DigitWidth = 1
	PageCount  = 1 << (DigitWidth * 4)
)

var options objekte_format.Options

func init() {
	options = objekte_format.Options{
		IncludeTai:           true,
		IncludeVerzeichnisse: true,
		PrintFinalSha:        true,
	}
}

type Store struct {
	erworben *konfig.Compiled
	path     string
	schnittstellen.VerzeichnisseFactory
	pages [PageCount]*Page
	ennui ennui.Ennui
}

type pageId struct {
	index uint8
	path  string
}

func MakeStore(
	s standort.Standort,
	k *konfig.Compiled,
	dir string,
	fff PageDelegateGetter,
) (i *Store, err error) {
	i = &Store{
		erworben:             k,
		path:                 dir,
		VerzeichnisseFactory: s,
	}

	if i.ennui, err = ennui.Make(s, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	for n := range i.pages {
		i.pages[n] = makePage(
			s.SansAge().SansCompression(),
			i.PageIdForIndex(uint8(n)),
			fff,
			i.ennui,
		)
	}

	return
}

func (i Store) PageIdForIndex(n uint8) (pid pageId) {
	pid.index = n
	pid.path = filepath.Join(i.path, fmt.Sprintf("%x", n))
	return
}

func (i *Store) ReadOne(k string, sk *sku.Transacted) (out *sku.Transacted, err error) {
	var loc ennui.Loc

	if loc, err = i.ennui.ReadOne(k, sk.GetMetadatei()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.readLoc(loc)
}

func (i *Store) readLoc(loc ennui.Loc) (sk *sku.Transacted, err error) {
	p := i.pages[loc.Page]

	var f *os.File

	if f, err = files.OpenFile(
		p.path,
		os.O_RDONLY,
		0o666,
	); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = f.Seek(int64(loc.Offset), io.SeekStart); err != nil {
		err = errors.Wrap(err)
		return
	}

	rb := catgut.MakeRingBuffer(f, 0)

	sk = sku.GetTransactedPool().Get()

	fo := objekte_format.Default()

	if _, err = fo.ParsePersistentMetadatei(
		rb,
		sk,
		options,
	); err != nil {
		err = errors.Wrapf(
			err,
			"Loc: %d, Readable: %s",
			loc.Offset,
			rb.PeekReadable().String()[:100],
		)

		return
	}

	return
}

// func (i *Store) ReadMany(string, *metadatei.Metadatei, *[]Loc) error {}
// func (i *Store) ReadAll(*metadatei.Metadatei, *[]Loc) error          {}

func (i Store) GetPage(n uint8) (p *Page, err error) {
	p = i.pages[n]
	return
}

func (i *Store) SetNeedsFlush() {
	for _, p := range i.pages {
		p.State = StateChanged
	}
}

func (i *Store) Flush() (err error) {
	errors.Log().Print("flushing")
	wg := iter.MakeErrorWaitGroup()

	wg.Do(i.ennui.Flush)

	for _, p := range i.pages {
		wg.Do(p.Flush)
	}

	return wg.GetError()
}

func (i *Store) AddVerzeichnisse(
	tz *sku.Transacted,
	v string,
) (err error) {
	var n uint8

	if n, err = i.PageForString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	var p *Page

	if p, err = i.GetPage(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	z := sku.GetTransactedPool().Get()

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

func (i *Store) ReadMany(
	ws ...schnittstellen.FuncIter[*sku.Transacted],
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

	w := pool.MakePooledChain[sku.Transacted](
		sku.GetTransactedPool(),
		ws...,
	)

	for n, p := range i.pages {
		wg.Add(1)

		go func(n int, p *Page, openFileCh chan struct{}) {
			defer wg.Done()
			defer func(c chan<- struct{}) {
				openFileCh <- struct{}{}
			}(openFileCh)

			for !isDone() {

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
