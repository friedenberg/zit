package stream_index

import (
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
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

var options object_inventory_format.Options

func init() {
	options = object_inventory_format.Options{
		Tai:           true,
		Verzeichnisse: true,
		PrintFinalSha: true,
	}
}

type Index struct {
	directoryLayout dir_layout.DirLayout
	mutable_config  *config.Compiled
	path            string
	interfaces.CacheIOFactory
	pages             [PageCount]Page
	historicalChanges []string
	probe_index
}

func MakeIndex(
	s dir_layout.DirLayout,
	k *config.Compiled,
	dir string,
) (i *Index, err error) {
	i = &Index{
		directoryLayout: s,
		mutable_config:  k,
		path:            dir,
		CacheIOFactory:  s,
	}

	if err = i.probe_index.Initialize(
		s,
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

func (i *Index) Initialize() (err error) {
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

func (i *Index) GetPage(n uint8) (p *Page) {
	p = &i.pages[n]
	return
}

func (i *Index) GetProbeIndex() *probe_index {
	return &i.probe_index
}

func (i *Index) SetNeedsFlushHistory(changes []string) {
	i.historicalChanges = changes
}

func (i *Index) Flush(
	printerHeader interfaces.FuncIter[string],
) (err error) {
	if len(i.historicalChanges) > 0 {
		if err = i.flushEverything(printerHeader); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = i.flushAdded(printerHeader); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *Index) flushAdded(
	printerHeader interfaces.FuncIter[string],
) (err error) {
	ui.Log().Print("flushing")
	wg := quiter.MakeErrorWaitGroupParallel()

	actualFlushCount := 0

	for n := range i.pages {
		if i.pages[n].hasChanges {
			ui.Log().Printf("actual flush for %d", n)
			actualFlushCount++
		}

		wg.Do(i.pages[n].MakeFlush(false))
	}

	if actualFlushCount > 0 {
		if err = printerHeader(
			fmt.Sprintf(
				"appending to index (%d/%d pages)",
				actualFlushCount,
				len(i.pages),
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	wg.DoAfter(i.Index.Flush)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if actualFlushCount > 0 {
		if err = printerHeader(
			fmt.Sprintf(
				"appended to index (%d/%d pages)",
				actualFlushCount,
				len(i.pages),
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *Index) flushEverything(
	printerHeader interfaces.FuncIter[string],
) (err error) {
	ui.Log().Print("flushing")
	wg := quiter.MakeErrorWaitGroupParallel()

	for n := range i.pages {
		wg.Do(i.pages[n].MakeFlush(true))
	}

	if err = printerHeader(
		fmt.Sprintf(
			"writing index (%d pages)",
			len(i.pages),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for n, change := range i.historicalChanges {
		if err = printerHeader(fmt.Sprintf("change: %s", change)); err != nil {
			err = errors.Wrap(err)
			return
		}

		if n == 99 {
			if err = printerHeader(
				fmt.Sprintf(
					"(%d more changes omitted)",
					len(i.historicalChanges)-100,
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			break
		}
	}

	wg.DoAfter(i.Index.Flush)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = printerHeader(
		fmt.Sprintf(
			"wrote index (%d pages)",
			len(i.pages),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Index) Add(
	z *sku.Transacted,
	v string,
	options sku.CommitOptions,
) (err error) {
	var n uint8

	if n, err = sha.PageIndexForString(DigitWidth, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.GetPage(n)

	if err = p.add(z, options); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Index) ReadOneSha(
	sh *sha.Sha,
	sk *sku.Transacted,
) (err error) {
	var loc object_probe_index.Loc

	if loc, err = s.readOneShaLoc(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.readOneLoc(loc, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Index) ReadManySha(
	sh *sha.Sha,
) (skus []*sku.Transacted, err error) {
	var locs []object_probe_index.Loc

	if locs, err = s.readManyShaLoc(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, loc := range locs {
		sk := sku.GetTransactedPool().Get()

		if err = s.readOneLoc(loc, sk); err != nil {
			err = errors.Wrapf(err, "Loc: %s", loc)
			return
		}

		skus = append(skus, sk)
	}

	return
}

func (s *Index) ObjectExists(
	id ids.IdLike,
) (err error) {
	var n uint8

	oid := id.String()

	if n, err = sha.PageIndexForString(DigitWidth, oid); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := s.GetPage(n)

	if _, ok := p.oids[oid]; ok {
		return
	}

	sh := sha.FromString(oid)
	defer sha.GetPool().Put(sh)

	if _, err = s.readOneShaLoc(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Index) ReadOneObjectId(
	id string,
	sk *sku.Transacted,
) (err error) {
	sh := sha.FromString(id)
	defer sha.GetPool().Put(sh)

	if err = s.ReadOneSha(sh, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Index) ReadManyObjectId(
	id string,
) (skus []*sku.Transacted, err error) {
	sh := sha.FromString(id)
	defer sha.GetPool().Put(sh)

	if skus, err = s.ReadManySha(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Index) ReadOneObjectIdTai(
	k interfaces.ObjectId,
	t ids.Tai,
) (sk *sku.Transacted, err error) {
	if t.IsEmpty() {
		err = collections.MakeErrNotFoundString(t.String())
		return
	}

	sh := sha.FromString(k.String() + t.String())
	defer sha.GetPool().Put(sh)

	sk = sku.GetTransactedPool().Get()

	if err = s.ReadOneSha(sh, sk); err != nil {
		err = errors.Wrapf(err, "ObjectId: %q, Tai: %q", k, t)
		return
	}

	return
}

func (s *Index) readOneLoc(
	loc object_probe_index.Loc,
	sk *sku.Transacted,
) (err error) {
	p := s.pages[loc.Page]

	if err = p.readOneRange(loc.Range, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Index) readFrom(
	qg sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
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

	w = pool.MakePooledChain(
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

					case quiter.IsStopIteration(err1):

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

func (i *Index) ReadQuery(
	qg sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return i.readFrom(qg, w)
}
