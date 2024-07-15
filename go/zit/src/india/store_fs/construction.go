package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeCwdFilesAll(
	k sku.Config,
	dp interfaces.FuncIter[*fd.FD],
	fileExtensions file_extensions.FileExtensions,
	st fs_home.Home,
	ofo object_inventory_format.Options,
) (fs *Store, err error) {
	fs = &Store{
		config:         k,
		deletedPrinter: dp,
		fs_home:        st,
		fileEncoder:    MakeFileEncoder(st, k),
		fileExtensions: fileExtensions,
		dir:            st.Cwd(),
		repos: collections_value.MakeMutableValueSet[*ObjectIdFDPair](
			nil,
		),
		types: collections_value.MakeMutableValueSet[*ObjectIdFDPair](nil),
		zettels: collections_value.MakeMutableValueSet[*ObjectIdFDPair](
			nil,
		),
		unsureZettels: collections_value.MakeMutableValueSet[*ObjectIdFDPair](
			nil,
		),
		tags: collections_value.MakeMutableValueSet[*ObjectIdFDPair](
			nil,
		),
		unsureBlobs: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		emptyDirectories: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		deleted: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		objectFormatOptions: ofo,
		metadataTextParser: object_metadata.MakeTextParser(
			st,
			nil, // TODO-P1 make akteFormatter
		),
	}

	if err = fs.readAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
