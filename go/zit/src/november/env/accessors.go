package env

import (
	"io"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
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
) interfaces.PoolValue[sku.ExternalLike] {
	if repoId.IsEmpty() {
		return nil
	}

	kid := repoId.GetRepoIdString()
	es, ok := u.externalStores[kid]

	if !ok {
		return pool.ManualPool[sku.ExternalLike]{
			FuncGet: func() sku.ExternalLike {
				return sku.GetTransactedPool().Get()
			},
			FuncPut: func(e sku.ExternalLike) {
				sku.GetTransactedPool().Put(e.(*sku.Transacted))
			},
		}
	} else {
		return es.GetExternalLikePool()
	}
}
