package object_probe_index

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
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
		AddMetadatei(*object_metadata.Metadatei, Loc) error
		ReadOneKey(string, *object_metadata.Metadatei) (Loc, error)
		ReadManyKeys(string, *object_metadata.Metadatei, *[]Loc) error
		ReadAll(*object_metadata.Metadatei, *[]Loc) error
		PrintAll() error
		errors.Flusher
	}
)

type Metadatei = object_metadata.Metadatei

const (
	DigitWidth = 1
	PageCount  = 1 << (DigitWidth * 4)
)

type object_probe_index struct {
	pages [PageCount]page
}

func MakePermitDuplicates(s fs_home.Standort, path string) (e *object_probe_index, err error) {
	e = &object_probe_index{}
	err = e.initialize(rowEqualerComplete{}, s, path)
	return
}

func MakeNoDuplicates(s fs_home.Standort, path string) (e *object_probe_index, err error) {
	e = &object_probe_index{}
	err = e.initialize(rowEqualerShaOnly{}, s, path)
	return
}

func (e *object_probe_index) initialize(
	equaler interfaces.Equaler1[*row],
	s fs_home.Standort,
	path string,
) (err error) {
	for i := range e.pages {
		p := &e.pages[i]
		p.initialize(equaler, s, sha.PageIdFromPath(uint8(i), path))
	}

	return
}

func (e *object_probe_index) GetEnnui() Ennui {
	return e
}

func (e *object_probe_index) AddMetadatei(m *Metadatei, loc Loc) (err error) {
	var shas map[string]*sha.Sha

	if shas, err = object_inventory_format.GetShasForMetadatei(m); err != nil {
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

func (e *object_probe_index) AddSha(sh *Sha, loc Loc) (err error) {
	return e.addSha(sh, loc)
}

func (e *object_probe_index) addSha(sh *Sha, loc Loc) (err error) {
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

func (e *object_probe_index) ReadOne(sh *Sha) (loc Loc, err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadOne(sh)
}

func (e *object_probe_index) ReadMany(sh *Sha, locs *[]Loc) (err error) {
	var i uint8

	if i, err = sha.PageIndexForSha(DigitWidth, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.pages[i].ReadMany(sh, locs)
}

func (e *object_probe_index) ReadOneKey(kf string, m *object_metadata.Metadatei) (loc Loc, err error) {
	var f object_inventory_format.FormatGeneric

	if f, err = object_inventory_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = object_inventory_format.GetShaForMetadatei(f, m); err != nil {
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

func (e *object_probe_index) ReadManyKeys(
	kf string,
	m *object_metadata.Metadatei,
	h *[]Loc,
) (err error) {
	var f object_inventory_format.FormatGeneric

	if f, err = object_inventory_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = object_inventory_format.GetShaForMetadatei(f, m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.ReadMany(sh, h)
}

func (e *object_probe_index) ReadAll(m *object_metadata.Metadatei, h *[]Loc) (err error) {
	var shas map[string]*sha.Sha

	if shas, err = object_inventory_format.GetShasForMetadatei(m); err != nil {
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

func (e *object_probe_index) PrintAll() (err error) {
	return
}

func (e *object_probe_index) Flush() (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	for i := range e.pages {
		p := &e.pages[i]
		wg.Do(p.Flush)
	}

	return wg.GetError()
}
