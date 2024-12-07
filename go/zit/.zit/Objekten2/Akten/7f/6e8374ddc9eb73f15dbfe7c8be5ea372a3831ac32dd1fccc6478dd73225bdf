package dir_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func (s DirLayout) ReadAllLevel2Files(
	p string,
	w interfaces.FuncIter[string],
) (err error) {
	if err = files.ReadDirNamesLevel2(
		files.MakeDirNameWriterIgnoringHidden(w),
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s DirLayout) ReadAllShas(
	p string,
	w interfaces.FuncIter[*sha.Sha],
) (err error) {
	wf := func(p string) (err error) {
		var sh *sha.Sha

		if sh, err = sha.MakeShaFromPath(p); err != nil {
			err = errors.Wrapf(err, "Path: %s", p)
			return
		}

		if err = w(sh); err != nil {
			err = errors.Wrapf(err, "Sha: %s", sh)
			return
		}

		return
	}

	if err = s.ReadAllLevel2Files(p, wf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s DirLayout) ReadAllShasForGenre(
	g interfaces.GenreGetter,
	w interfaces.FuncIter[*sha.Sha],
) (err error) {
	var p string

	if p, err = s.DirObjectGenre(g); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadAllShas(p, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
