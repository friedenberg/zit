package remote_http

import (
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Client struct {
	http.Client
	Repo *local_working_copy.Repo
	// *local_working_copy.Repo
}

func (repo *Client) GetRepoType() repo_type.Type {
	return repo_type.TypeUnknown
}

func (repo *Client) GetInventoryListStore() sku.InventoryListStore {
	return nil
	// return repo
}

func (repo *Client) GetBlobStore() interfaces.BlobStore {
	return &HTTPBlobStore{repo: repo}
}
