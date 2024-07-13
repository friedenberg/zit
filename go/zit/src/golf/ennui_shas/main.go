package ennui_shas

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
)

type (
	Sha = sha.Sha

	commonInterface interface {
		AddSha(*sha.Sha, *sha.Sha) error
		ReadOne(left *Sha) (right *sha.Sha, err error)
		ReadMany(left *Sha, rights *[]*sha.Sha) (err error)
	}

	pageInterface interface {
		GetEnnuiPage() pageInterface
		commonInterface
		PrintAll() error
		errors.Flusher
	}

	Ennui interface {
		GetEnnui() Ennui
		commonInterface
		PrintAll() error
		errors.Flusher
	}
)

const (
	DigitWidth = 1
	PageCount  = 1 << (DigitWidth * 4)
)

type ennui struct {
	rowSize int
	pages   [PageCount]page
}

func MakePermitDuplicates(
	s fs_home.Standort,
	path string,
) (e *ennui, err error) {
	e = &ennui{}
	e.rowSize = RowSize
	err = e.initialize(rowEqualerComplete{}, s, path)
	return
}

func MakeNoDuplicates(s fs_home.Standort, path string) (e *ennui, err error) {
	e = &ennui{}
	e.rowSize = RowSize
	err = e.initialize(rowEqualerShaOnly{}, s, path)
	return
}

func (e *ennui) initialize(
	equaler interfaces.Equaler1[*row],
	s fs_home.Standort,
	path string,
) (err error) {
	for i := range e.pages {
		p := &e.pages[i]
		p.initialize(equaler, s, sha.PageIdFromPath(uint8(i), path), e.rowSize)
	}

	return
}

func (e *ennui) GetEnnui() Ennui {
	return e
}

func (e *ennui) AddSha(left, right *Sha) (err error) {
	return e.addSha(left, right)
}

func (e *ennui) addSha(left, right *Sha) (err error) {
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

func (e *ennui) ReadOne(left *Sha) (right *sha.Sha, err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, left); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadOne(left)
}

func (e *ennui) ReadMany(sh *Sha, locs *[]*sha.Sha) (err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadMany(sh, locs)
}

func (e *ennui) PrintAll() (err error) {
	return
}

func (e *ennui) Flush() (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	for i := range e.pages {
		p := &e.pages[i]
		wg.Do(p.Flush)
	}

	return wg.GetError()
}
