package ennui_shas

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
)

type (
	Sha = sha.Sha

	commonInterface interface {
		AddSha(*sha.Sha, *Loc) error
		ReadOne(sh *Sha) (loc *Loc, err error)
		ReadMany(sh *Sha, locs *[]*Loc) (err error)
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
	pages [PageCount]page
}

func MakePermitDuplicates(
	s standort.Standort,
	path string,
) (e *ennui, err error) {
	e = &ennui{}
	err = e.initialize(rowEqualerComplete{}, s, path)
	return
}

func MakeNoDuplicates(s standort.Standort, path string) (e *ennui, err error) {
	e = &ennui{}
	err = e.initialize(rowEqualerShaOnly{}, s, path)
	return
}

func (e *ennui) initialize(
	equaler schnittstellen.Equaler1[*row],
	s standort.Standort,
	path string,
) (err error) {
	for i := range e.pages {
		p := &e.pages[i]
		p.initialize(equaler, s, sha.PageIdFromPath(uint8(i), path))
	}

	return
}

func (e *ennui) GetEnnui() Ennui {
	return e
}

func (e *ennui) AddSha(sh *Sha, loc *Loc) (err error) {
	return e.addSha(sh, loc)
}

func (e *ennui) addSha(sh *Sha, loc *Loc) (err error) {
	if sh.IsNull() {
		return
	}

	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].AddSha(sh, loc)
}

func (e *ennui) ReadOne(sh *Sha) (loc *Loc, err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadOne(sh)
}

func (e *ennui) ReadMany(sh *Sha, locs *[]*Loc) (err error) {
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
