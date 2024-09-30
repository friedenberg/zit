package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) Open(
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()

	if m.IncludesMetadata() {
		wg.Do(func() error {
			return s.openZettels(ph, zsc)
		})
	}

	if m.IncludesBlob() {
		wg.Do(func() error {
			return s.openBlob(ph, zsc)
		})
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) openZettels(
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	var filesZettels []string

	if filesZettels, err = s.ToSliceFilesZettelen(zsc); err != nil {
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

	if err = openVimOp.Run(ph, filesZettels...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) openBlob(
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	var filesBlobs []string

	if filesBlobs, err = s.ToSliceFilesBlobs(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	opOpenFiles := OpenFiles{}

	if err = opOpenFiles.Run(ph, filesBlobs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
