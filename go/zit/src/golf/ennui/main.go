package ennui

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
)

type (
	Sha = sha.Sha

	commonInterface interface {
		AddSha(*sha.Sha, Loc) error
		ReadOne(sh *Sha) (loc Loc, err error)
		ReadMany(sh *Sha, locs *[]Loc) (err error)
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
		AddMetadatei(*metadatei.Metadatei, Loc) error
		ReadOneKey(string, *metadatei.Metadatei) (Loc, error)
		ReadManyKeys(string, *metadatei.Metadatei, *[]Loc) error
		ReadAll(*metadatei.Metadatei, *[]Loc) error
		PrintAll() error
		errors.Flusher
	}
)

type Metadatei = metadatei.Metadatei

const (
	DigitWidth = 1
	PageCount  = 1 << (DigitWidth * 4)
)

type ennui struct {
	pages [PageCount]page
}

func MakePermitDuplicates(s standort.Standort, path string) (e *ennui, err error) {
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

func (e *ennui) AddMetadatei(m *Metadatei, loc Loc) (err error) {
	var shas map[string]*sha.Sha

	if shas, err = objekte_format.GetShasForMetadatei(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, s := range shas {
		if err = e.addSha(s, loc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *ennui) AddSha(sh *Sha, loc Loc) (err error) {
	return e.addSha(sh, loc)
}

func (e *ennui) addSha(sh *Sha, loc Loc) (err error) {
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

func (e *ennui) ReadOne(sh *Sha) (loc Loc, err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadOne(sh)
}

func (e *ennui) ReadMany(sh *Sha, locs *[]Loc) (err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadMany(sh, locs)
}

func (e *ennui) ReadOneKey(kf string, m *metadatei.Metadatei) (loc Loc, err error) {
	var f objekte_format.FormatGeneric

	if f, err = objekte_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = objekte_format.GetShaForMetadatei(f, m); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer sha.GetPool().Put(sh)

	if loc, err = e.ReadOne(sh); err != nil {
		err = errors.Wrapf(err, "Key: %s", kf)
		return
	}

	return
}

func (e *ennui) ReadManyKeys(
	kf string,
	m *metadatei.Metadatei,
	h *[]Loc,
) (err error) {
	var f objekte_format.FormatGeneric

	if f, err = objekte_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = objekte_format.GetShaForMetadatei(f, m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.ReadMany(sh, h)
}

func (e *ennui) ReadAll(m *metadatei.Metadatei, h *[]Loc) (err error) {
	var shas map[string]*sha.Sha

	if shas, err = objekte_format.GetShasForMetadatei(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()

	for k, s := range shas {
		s := s
		wg.Do(
			func() (err error) {
				var loc Loc

				if loc, err = e.ReadOne(s); err != nil {
					err = errors.Wrapf(err, "Key: %s", k)
					return
				}

				*h = append(*h, loc)

				return
			},
		)
	}

	return wg.GetError()
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
