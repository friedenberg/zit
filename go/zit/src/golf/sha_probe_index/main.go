package sha_probe_index

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

type (
	Sha = sha.Sha

	commonInterface interface {
		AddSha(*sha.Sha, *sha.Sha) error
		ReadOne(left *Sha) (right *sha.Sha, err error)
		ReadMany(left *Sha, rights *[]*sha.Sha) (err error)
	}

	pageInterface interface {
		GetIndexPage() pageInterface
		commonInterface
		PrintAll() error
		errors.Flusher
	}

	Index interface {
		GetIndex() Index
		commonInterface
		PrintAll() error
		errors.Flusher
	}
)

const (
	DigitWidth = 1
	PageCount  = 1 << (DigitWidth * 4)
)

type object_probe_index struct {
	rowSize int
	pages   [PageCount]page
}

func MakePermitDuplicates(
	s dir_layout.DirLayout,
	path string,
) (e *object_probe_index, err error) {
	e = &object_probe_index{}
	e.rowSize = RowSize
	err = e.initialize(rowEqualerComplete{}, s, path)
	return
}

func MakeNoDuplicates(s dir_layout.DirLayout, path string) (e *object_probe_index, err error) {
	e = &object_probe_index{}
	e.rowSize = RowSize
	err = e.initialize(rowEqualerShaOnly{}, s, path)
	return
}

func (e *object_probe_index) initialize(
	equaler interfaces.Equaler1[*row],
	s dir_layout.DirLayout,
	path string,
) (err error) {
	for i := range e.pages {
		p := &e.pages[i]
		p.initialize(equaler, s, sha.PageIdFromPath(uint8(i), path), e.rowSize)
	}

	return
}

func (e *object_probe_index) GetIndex() Index {
	return e
}

func (e *object_probe_index) AddSha(left, right *Sha) (err error) {
	return e.addSha(left, right)
}

func (e *object_probe_index) addSha(left, right *Sha) (err error) {
	if left.IsNull() {
		return
	}

	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, left); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].AddSha(left, right)
}

func (e *object_probe_index) ReadOne(left *Sha) (right *sha.Sha, err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, left); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadOne(left)
}

func (e *object_probe_index) ReadMany(sh *Sha, locs *[]*sha.Sha) (err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadMany(sh, locs)
}

func (e *object_probe_index) PrintAll() (err error) {
	return
}

func (e *object_probe_index) Flush() (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()

	for i := range e.pages {
		p := &e.pages[i]
		wg.Do(p.Flush)
	}

	return wg.GetError()
}
