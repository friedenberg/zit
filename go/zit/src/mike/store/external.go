package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

type ObjekteOptions = sku.ObjekteOptions

func (s *Store) ReadOneKennungExternal(
	o ObjekteOptions,
	k1 schnittstellen.StringerGattungKastenGetter,
	sk *sku.Transacted,
) (el sku.ExternalLike, err error) {
	switch k1.GetKasten().GetKastenString() {
	case "chrome":
		// TODO populate with chrome kasten
		ui.Debug().Print("would populate from chrome")

	default:
		if el, err = s.cwdFiles.ReadKennung(o, k1, sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) Open(
	kasten schnittstellen.KastenGetter,
	m checkout_mode.Mode,
	ph schnittstellen.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	switch kasten.GetKasten().GetKastenString() {
	case "chrome":
		err = todo.Implement()

	default:
		if err = s.OpenFS(m, ph, zsc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) OpenFS(
	m checkout_mode.Mode,
	ph schnittstellen.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	var filesZettelen []string

	if filesZettelen, err = store_fs.ToSliceFilesZettelen(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := store_fs.OpenVim{
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
