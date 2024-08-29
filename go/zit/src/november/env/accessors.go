package env

import (
	"io"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Env) GetTime() time.Time {
	return time.Now()
}

func (u *Env) GetConfig() *config.Compiled {
	return &u.config
}

func (u *Env) GetDormantIndex() *dormant_index.Index {
	return &u.dormantIndex
}

func (u *Env) In() io.Reader {
	return u.in
}

func (u *Env) Out() interfaces.WriterAndStringWriter {
	return u.out
}

func (u *Env) Err() interfaces.WriterAndStringWriter {
	return u.err
}

func (u *Env) GetFSHome() fs_home.Home {
	return u.fs_home
}

func (u *Env) GetStore() *store.Store {
	return &u.store
}

func (u *Env) GetExternalLikePoolForRepoId(
	repoId ids.RepoId,
) (of external_store.ObjectFactory) {
	if repoId.IsEmpty() {
		return
	}

	kid := repoId.GetRepoIdString()
	es, ok := u.externalStores[kid]

	if ok {
		of.PoolValue = es.GetExternalLikePool()
		of.Resetter3 = es.GetExternalLikeResetter3()
	}

	return
}
