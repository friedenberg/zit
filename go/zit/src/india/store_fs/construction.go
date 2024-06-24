package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
)

func makeCwdFiles(
	fileExtensions file_extensions.FileExtensions,
	st standort.Standort,
) (fs *Store) {
	fs = &Store{
		standort:       st,
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
		metadateiTextParser: metadatei.MakeTextParser(
			st,
			nil, // TODO-P1 make akteFormatter
		),
	}

	return
}

func MakeCwdFilesAll(
	fileExtensions file_extensions.FileExtensions,
	st standort.Standort,
) (fs *Store, err error) {
	fs = makeCwdFiles(fileExtensions, st)
	err = fs.readAll()
	return
}

func MakeCwdFilesExactly(
	fileExtensions file_extensions.FileExtensions,
	st standort.Standort,
	files ...string,
) (fs *Store, err error) {
	fs = makeCwdFiles(fileExtensions, st)
	err = fs.readInputFiles(files...)
	return
}
