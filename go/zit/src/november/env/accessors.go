package env

import (
	"io"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Local) GetTime() time.Time {
	return time.Now()
}

func (u *Local) GetConfig() *config.Compiled {
	return &u.config
}

func (u *Local) GetDormantIndex() *dormant_index.Index {
	return &u.dormantIndex
}

func (u *Local) In() io.Reader {
	return u.in.File
}

func (u *Local) Out() interfaces.WriterAndStringWriter {
	return u.out.File
}

func (u *Local) Err() interfaces.WriterAndStringWriter {
	return u.err.File
}

func (u *Local) GetDirectoryLayout() dir_layout.DirLayout {
	return u.dirLayout
}

func (u *Local) GetDirLayoutPrimitive() dir_layout.Primitive {
	return u.dirLayoutPrimitive
}

func (u *Local) GetStore() *store.Store {
	return &u.store
}

func (u *Local) GetExternalLikePoolForRepoId(
	repoId ids.RepoId,
) (of sku.ObjectFactory) {
	return
}

func (u *Local) GetFileEncoder() store_fs.FileEncoder {
	return u.fileEncoder
}
