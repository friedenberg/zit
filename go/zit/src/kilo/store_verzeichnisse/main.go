package store_verzeichnisse

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/juliett/query"
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
		Tai:           true,
		Verzeichnisse: true,
		PrintFinalSha: true,
	}
}

type Store struct {
	standort standort.Standort
	erworben *konfig.Compiled
	path     string
	schnittstellen.VerzeichnisseFactory
	pages [PageCount]Page
	ennuiStore
}

func MakeStore(
	s standort.Standort,
	k *konfig.Compiled,
	dir string,
	persistentMetadateiFormat objekte_format.Format,
	options objekte_format.Options,
) (i *Store, err error) {
	i = &Store{
		standort:             s,
		erworben:             k,
		path:                 dir,
		VerzeichnisseFactory: s,
	}

	if err = i.ennuiStore.Initialize(
		s,
		persistentMetadateiFormat,
		options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Store) Initialize() (err error) {
	for n := range i.pages {
		i.pages[n].initialize(
			PageId{
				Prefix: "Page",
				Dir:    i.path,
				Index:  uint8(n),
			},
			i,
		)
	}

	return
}

func (s *Store) applyKonfig(z *sku.Transacted) error {
	if !s.erworben.HasChanges() {
		return nil
	}

	return s.erworben.ApplyToSku(z)
}

// func (i *Store) ReadMany(string, *metadatei.Metadatei, *[]Loc) error {}
// func (i *Store) ReadAll(*metadatei.Metadatei, *[]Loc) error          {}

func (i *Store) GetPagePair(n uint8) (p *Page) {
	p = &i.pages[n]
	return
}

func (i *Store) SetNeedsFlushHistory() {
	for n := range i.pages {
		i.pages[n].SetNeedsFlushHistory()
	}
}

func (i *Store) Flush(
	printerHeader schnittstellen.FuncIter[string],
) (err error) {
	errors.Log().Print("flushing")
	wg := iter.MakeErrorWaitGroupParallel()

	actualFlush := false

	for n := range i.pages {
		if i.pages[n].hasChanges {
			log.Log().Printf("actual flush for %d", n)
			actualFlush = true
		}

		wg.Do(i.pages[n].Flush)
	}

	if actualFlush {
		if err = printerHeader("writing index"); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	wg.DoAfter(i.kennung.Flush)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if actualFlush {
		if err = printerHeader("wrote index"); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *Store) Add(
	z *sku.Transacted,
	v string,
	mode objekte_mode.Mode,
) (err error) {
	var n uint8

	if n, err = sha.PageIndexForString(DigitWidth, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.GetPagePair(n)

	if err = p.add(z, mode); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Store) readFrom(
	qg *query.Group,
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
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

	w = pool.MakePooledChain[sku.Transacted](
		sku.GetTransactedPool(),
		w,
	)

	for n := range i.pages {
		wg.Add(1)

		go func(p *Page, openFileCh chan struct{}) {
			defer wg.Done()
			defer func() {
				openFileCh <- struct{}{}
			}()

			for !isDone() {

				var err1 error

				if err1 = p.CopyJustHistory(
					qg,
					w,
				); err1 != nil {
					if isDone() {
						break
					}

					switch {
					case errors.IsTooManyOpenFiles(err1):
						<-openFileCh
						continue

					case iter.IsStopIteration(err1):

					default:
						me.Add(err1)
					}
				}

				break
			}
		}(&i.pages[n], ch)
	}

	wg.Wait()

	if me.Len() > 0 {
		err = me
	}

	return
}

func (i *Store) ReadQuery(
	qg *query.Group,
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return i.readFrom(qg, w)
}
