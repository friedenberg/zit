package repo_remote

import (
	"bufio"
	"net/http"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

// type HTTPClient interface {
// 	Do(req *http.Request) (*http.Response, error)
// }

type HTTP struct {
	http.Client
	remote *repo_local.Repo
}

func (remote *HTTP) GetRepo() repo.Repo {
	return remote
}

func (remote *HTTP) GetBlobStore() interfaces.BlobStore {
	return &HTTPBlobStore{remote: remote}
}

func (remote *HTTP) MakeQueryGroup(
	metaBuilder any,
	repoId ids.RepoId,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Group, err error) {
	err = todo.Implement()
	return
}

func (remote *HTTP) MakeInventoryList(
	qg *query.Group,
) (list *sku.List, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		remote.remote.Context,
		"GET",
		"/inventory_list",
		strings.NewReader(qg.String()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = remote.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	bf := remote.remote.GetStore().GetInventoryListStore().FormatForVersion(
		remote.remote.GetConfig().GetStoreVersion(),
	)

	list = sku.MakeList()

	if err = inventory_list_blobs.ReadInventoryListBlob(
		bf,
		bufio.NewReader(response.Body),
		list,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (remoteHTTP *HTTP) PullQueryGroupFromRemote(
	remote repo.Repo,
	qg *query.Group,
	printCopies bool,
) (err error) {
	err = todo.Implement()
	return
}

func (remote *HTTP) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	err = todo.Implement()
	return
}

type HTTPBlobStore struct {
	remote *HTTP
}

func (blobStore *HTTPBlobStore) GetBlobStore() interfaces.BlobStore {
	return blobStore
}

func (blobStore *HTTPBlobStore) HasBlob(sh interfaces.Sha) (ok bool) {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			blobStore.remote.remote.Context,
			"HEAD",
			"/blobs",
			strings.NewReader(sh.GetShaLike().GetShaString()),
		); err != nil {
			blobStore.remote.remote.CancelWithError(err)
		}
	}

	var response *http.Response

	{
		var err error

		if response, err = blobStore.remote.Do(request); err != nil {
			blobStore.remote.remote.CancelWithError(err)
		}
	}

	ok = response.StatusCode == http.StatusNoContent

	return
}

func (blobStore *HTTPBlobStore) BlobWriter() (w interfaces.ShaWriteCloser, err error) {
	err = todo.Implement()
	return
}

func (blobStore *HTTPBlobStore) BlobReader(
	sh interfaces.Sha,
) (r interfaces.ShaReadCloser, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		blobStore.remote.remote.Context,
		"GET",
		"/blobs",
		strings.NewReader(sh.GetShaLike().GetShaString()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = blobStore.remote.Do(request); err != nil {
		err = errors.Wrap(err)
		return
	}

	if response.StatusCode != http.StatusOK {
		err = errors.Errorf("remote error: %d", response.StatusCode)
		return
	}

	r = sha.MakeReadCloser(response.Body)

	return
}
