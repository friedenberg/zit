package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func makeCwdFiles(
	k sku.Konfig,
	storeFuncs sku.StoreFuncs,
	fileExtensions file_extensions.FileExtensions,
	st standort.Standort,
	ofo objekte_format.Options,
) (fs *Store) {
	fs = &Store{
		konfig:         k,
		storeFuncs:     storeFuncs,
		standort:       st,
		fileEncoder:    MakeFileEncoder(st, k),
		fileExtensions: fileExtensions,
		dir:            st.Cwd(),
		kisten: collections_value.MakeMutableValueSet[*KennungFDPair](
			nil,
		),
		typen: collections_value.MakeMutableValueSet[*KennungFDPair](nil),
		zettelen: collections_value.MakeMutableValueSet[*KennungFDPair](
			nil,
		),
		unsureZettelen: collections_value.MakeMutableValueSet[*KennungFDPair](
			nil,
		),
		etiketten: collections_value.MakeMutableValueSet[*KennungFDPair](
			nil,
		),
		unsureAkten: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		emptyDirectories: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		deleted: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		objekteFormatOptions: ofo,
		metadateiTextParser: metadatei.MakeTextParser(
			st,
			nil, // TODO-P1 make akteFormatter
		),
	}

	return
}

func MakeCwdFilesAll(
	k sku.Konfig,
	storeFuncs sku.StoreFuncs,
	fileExtensions file_extensions.FileExtensions,
	st standort.Standort,
	ofo objekte_format.Options,
) (fs *Store, err error) {
	fs = makeCwdFiles(k, storeFuncs, fileExtensions, st, ofo)
	err = fs.readAll()
	return
}
