package object_probe_index

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
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
		GetObjectProbeIndexPage() pageInterface
		commonInterface
		PrintAll() error
		errors.Flusher
	}

	Index interface {
		GetObjectProbeIndex() Index
		commonInterface
		PrintAll() error
		errors.Flusher
	}
)

type Metadata = object_metadata.Metadata

const (
	DigitWidth = 1
	PageCount  = 1 << (DigitWidth * 4)
)

type object_probe_index struct {
	pages [PageCount]page
}

func MakePermitDuplicates(s dir_layout.DirLayout, path string) (e *object_probe_index, err error) {
	e = &object_probe_index{}
	err = e.initialize(rowEqualerComplete{}, s, path)
	return
}

func MakeNoDuplicates(s dir_layout.DirLayout, path string) (e *object_probe_index, err error) {
	e = &object_probe_index{}
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
		p.initialize(equaler, s, sha.PageIdFromPath(uint8(i), path))
	}

	return
}

func (e *object_probe_index) GetObjectProbeIndex() Index {
	return e
}

func (e *object_probe_index) AddMetadata(m *Metadata, loc Loc) (err error) {
	var shas map[string]*sha.Sha

	if shas, err = object_inventory_format.GetShasForMetadata(m); err != nil {
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

func (e *object_probe_index) ReadOneKey(kf string, m *object_metadata.Metadata) (loc Loc, err error) {
	var f object_inventory_format.FormatGeneric

	if f, err = object_inventory_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = object_inventory_format.GetShaForMetadata(f, m); err != nil {
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
	m *object_metadata.Metadata,
	h *[]Loc,
) (err error) {
	var f object_inventory_format.FormatGeneric

	if f, err = object_inventory_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = object_inventory_format.GetShaForMetadata(f, m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return e.ReadMany(sh, h)
}

func (e *object_probe_index) ReadAll(m *object_metadata.Metadata, h *[]Loc) (err error) {
	var shas map[string]*sha.Sha

	if shas, err = object_inventory_format.GetShasForMetadata(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := quiter.MakeErrorWaitGroupParallel()

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
	for i := range e.pages {
		p := &e.pages[i]

		if err = p.PrintAll(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

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
