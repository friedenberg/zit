package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeCwdFilesAll(
	k sku.Config,
	dp interfaces.FuncIter[*fd.FD],
	fileExtensions file_extensions.FileExtensions,
	st dir_layout.DirLayout,
	ofo object_inventory_format.Options,
) (fs *Store, err error) {
	fs = &Store{
		config:         k,
		deletedPrinter: dp,
		dirLayout:        st,
		fileEncoder:    MakeFileEncoder(st, k),
		fileExtensions: fileExtensions,
		dir:            st.Cwd(),
		dirItems:       makeObjectsWithDir(st.Cwd(), fileExtensions, st),
		deleted: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		objectFormatOptions: ofo,
		metadataTextParser: object_metadata.MakeTextParser(
			st,
			nil,
		),
	}

	return
}
