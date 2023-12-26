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
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/golf/kennung_index"
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
	standort standort.Standort
	erworben *konfig.Compiled
	path     string
	schnittstellen.VerzeichnisseFactory
	pages [PageCount]PageTuple
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
	ki kennung_index.Index,
) (i *Store, err error) {
	i = &Store{
		standort:             s,
		erworben:             k,
		path:                 dir,
		VerzeichnisseFactory: s,
	}

	if i.ennui, err = ennui.Make(s, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	for n := range i.pages {
		i.pages[n].initialize(uint8(n), i, ki)
	}

	return
}

func (s *Store) applyKonfig(z *sku.Transacted) error {
	if !s.erworben.HasChanges() {
		return nil
	}

	return s.erworben.ApplyToSku(z)
}

func (i *Store) PageIdForIndex(n uint8, isSchwanz bool) (pid pageId) {
	pid.index = n

	if isSchwanz {
		pid.path = filepath.Join(i.path, fmt.Sprintf("Schwanz-%x", n))
	} else {
		pid.path = filepath.Join(i.path, fmt.Sprintf("All-%x", n))
	}

	return
}

func (i *Store) GetEnnui() ennui.Ennui {
  return i.ennui
}

func (i *Store) ReadOne(sh *sha.Sha) (out *sku.Transacted, err error) {
	var loc ennui.Loc

	if loc, err = i.ennui.ReadOneSha(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.readLoc(loc)
}

func (i *Store) ReadOneKey(k string, sk *sku.Transacted) (out *sku.Transacted, err error) {
	var loc ennui.Loc

	if loc, err = i.ennui.ReadOneKey(k, sk.GetMetadatei()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.readLoc(loc)
}

func (i *Store) readLoc(loc ennui.Loc) (sk *sku.Transacted, err error) {
	p := &i.pages[loc.Page].All

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

func (i *Store) GetPagePair(n uint8) (p *PageTuple) {
	p = &i.pages[n]
	return
}

func (i *Store) SetNeedsFlush() {
	for n := range i.pages {
		i.pages[n].All.State = StateChanged
		i.pages[n].Schwanzen.State = StateChanged
	}
}

func (i *Store) Flush() (err error) {
	errors.Log().Print("flushing")
	wg := iter.MakeErrorWaitGroup()

	wg.Do(i.ennui.Flush)

	for n := range i.pages {
		wg.Do(i.pages[n].Flush)
	}

	return wg.GetError()
}

func (i *Store) Add(
	z *sku.Transacted,
	v string,
) (err error) {
	var n uint8

	if n, err = i.PageForString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.GetPagePair(n)

	if err = p.All.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = p.Schwanzen.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Store) readFrom(
	s kennung.Sigil,
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

	for n := range i.pages {
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
		}(n, i.pages[n].PageForSigil(s), ch)
	}

	wg.Wait()

	if me.Len() > 0 {
		err = me
	}

	return
}

func (i *Store) ReadSchwanzen(
	ws ...schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return i.readFrom(kennung.SigilSchwanzen, ws...)
}

func (i *Store) ReadAll(
	ws ...schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return i.readFrom(kennung.SigilHistory, ws...)
}
