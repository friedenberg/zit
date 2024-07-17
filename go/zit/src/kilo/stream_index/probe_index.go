package stream_index

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type probe_index struct {
	fs_home fs_home.Home
	object_probe_index.Index
}

func (s *probe_index) Initialize(
	fs_home fs_home.Home,
) (err error) {
	s.fs_home = fs_home

	if s.Index, err = object_probe_index.MakeNoDuplicates(
		s.fs_home,
		s.fs_home.DirVerzeichnisseVerweise(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *probe_index) Flush() (err error) {
	if err = s.Index.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *probe_index) readOneShaLoc(
	sh *sha.Sha,
) (loc object_probe_index.Loc, err error) {
	if loc, err = s.Index.ReadOne(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *probe_index) saveOneLoc(
	o *sku.Transacted,
	loc object_probe_index.Loc,
) (err error) {
	if err = s.saveOneLocString(
		o,
		o.GetObjectId().String(),
		loc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.saveOneLocString(
		o,
		o.GetObjectId().String()+o.GetTai().String(),
		loc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *probe_index) saveOneLocString(
	o *sku.Transacted,
	str string,
	loc object_probe_index.Loc,
) (err error) {
	sh := sha.FromString(str)
	defer sha.GetPool().Put(sh)

	ui.Log().Print(str, sh, o, loc)

	if err = s.Index.AddSha(sh, loc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
