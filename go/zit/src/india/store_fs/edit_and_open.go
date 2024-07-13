package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) Open(
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	if m.IncludesObjekte() {
		wg.Do(func() error {
			return s.openZettelen(ph, zsc)
		})
	}

	if m.IncludesAkte() {
		wg.Do(func() error {
			return s.openAkten(ph, zsc)
		})
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) openZettelen(
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	var filesZettelen []string

	if filesZettelen, err = ToSliceFilesZettelen(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := Open{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-zettel").
			WithInsertMode().
			Build(),
	}

	if err = openVimOp.Run(ph, filesZettelen...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) openAkten(
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	var filesAkten []string

	if filesAkten, err = ToSliceFilesAkten(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	opOpenFiles := OpenFiles{}

	if err = opOpenFiles.Run(ph, filesAkten...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
