package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ExternalInfo struct {
	ExternalType ids.Type

	// TODO add support for querying the below
	ids.RepoId
	external_state.State
	ExternalObjectId ids.ObjectId
}

func (t *ExternalInfo) GetRepoId() ids.RepoId {
	return t.RepoId
}

func (t *ExternalInfo) GetExternalObjectId() ids.ExternalObjectId {
	return &t.ExternalObjectId
}

func (t *ExternalInfo) GetExternalState() external_state.State {
	return external_state.Unknown
}
