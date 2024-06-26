package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) Open(
	m checkout_mode.Mode,
	ph schnittstellen.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	var filesZettelen []string

	if filesZettelen, err = ToSliceFilesZettelen(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := OpenVim{
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
