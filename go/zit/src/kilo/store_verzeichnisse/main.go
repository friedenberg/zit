package store_verzeichnisse

import (
	"os"
	"path"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/src/golf/kennung_index"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/juliett/konfig"
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
	pages                   [PageCount]PageTuple
	tomlPages               [PageCount]TomlPageTuple
	ennuiShas, ennuiKennung ennui.Ennui
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

	if err = i.Initialize(ki); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Store) Initialize(ki kennung_index.Index) (err error) {
	if i.ennuiShas, err = ennui.MakePermitDuplicates(
		i.standort,
		path.Join(i.path, "EnnuiShas"),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if i.ennuiKennung, err = ennui.MakeNoDuplicates(
		i.standort,
		path.Join(i.path, "EnnuiKennung"),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for n := range i.pages {
		i.pages[n].initialize(
			PageId{
				Prefix: "Page",
				Dir:    i.path,
				Index:  uint8(n),
			},
			i,
			ki,
		)

		i.tomlPages[n].initialize(
			PageId{
				Prefix: "Page",
				Dir:    i.path,
				Index:  uint8(n),
			},
			i,
			ki,
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

func (i *Store) GetEnnuiShas() ennui.Ennui {
	return i.ennuiShas
}

func (i *Store) GetEnnuiKennung() ennui.Ennui {
	return i.ennuiKennung
}

func (i *Store) ExistsOneSha(sh *sha.Sha) (err error) {
	if _, err = i.ennuiShas.ReadOne(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	err = collections.ErrExists

	return
}

func (i *Store) ReadOneShas(sh *sha.Sha) (out *sku.Transacted, err error) {
	var loc ennui.Loc

	if loc, err = i.ennuiShas.ReadOne(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.readLoc(loc)
}

func (i *Store) ReadOneKennung(
	h schnittstellen.Stringer,
) (out *sku.Transacted, err error) {
	sh := sha.FromString(h.String())
	defer sha.GetPool().Put(sh)

	var loc ennui.Loc

	if loc, err = i.ennuiKennung.ReadOne(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.readLoc(loc)
}

func (i *Store) ReadOneAll(
	mg metadatei.Getter,
	kennungPtr kennung.Kennung,
) (out []ennui.Loc, err error) {
	var locKennung ennui.Loc

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(func() (err error) {
		sh := sha.FromString(kennungPtr.String())
		defer sha.GetPool().Put(sh)

		if locKennung, err = i.ennuiKennung.ReadOne(sh); err != nil {
			if collections.IsErrNotFound(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		return
	})

	wg.Do(func() (err error) {
		if err = i.ennuiShas.ReadAll(mg.GetMetadatei(), &out); err != nil {
			if collections.IsErrNotFound(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		return
	})

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !locKennung.IsEmpty() {
		out = append(out, locKennung)
	}

	return
}

func (i *Store) readLoc(loc ennui.Loc) (sk *sku.Transacted, err error) {
	p := &i.pages[loc.Page]

	var f *os.File

	if f, err = files.OpenFile(
		p.Path(),
		os.O_RDONLY,
		0o666,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	coder := sku_fmt.Binary{QueryGroup: &sku_fmt.NopSigil{Sigil: kennung.SigilAll}}

	sk = sku.GetTransactedPool().Get()

	if _, err = coder.ReadFormatExactly(
		f,
		loc,
		&sku_fmt.Sku{Transacted: sk},
	); err != nil {
		sku.GetTransactedPool().Put(sk)
		sk = nil
		err = errors.Wrapf(err, "%s", loc)
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
			actualFlush = true
		}

		wg.Do(i.pages[n].Flush)
	}

	wg.DoAfter(i.ennuiShas.Flush)
	wg.DoAfter(i.ennuiKennung.Flush)

	if actualFlush {
		if err = printerHeader("writing index"); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

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

		go func(p *PageTuple, openFileCh chan struct{}) {
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
