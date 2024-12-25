package repo_remote

import (
	"bufio"
	"net"
	"net/http"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

func MakeRemoteHTTPFromXDGDotenvPath(
	context errors.Context,
	config *config.Compiled,
	xdgDotenvPath string,
) (remoteHTTP *HTTP, err error) {
	var remote *repo_local.Repo

	if remote, err = repo_local.MakeFromConfigAndXDGDotenvPath(
		context,
		config,
		xdgDotenvPath,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	remoteHTTP = &HTTP{
		remote: remote,
	}

	if remoteHTTP.unixSocket, err = remote.InitializeUnixSocket(
		net.ListenConfig{},
		"",
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	go func() {
		if err := remote.Serve(remoteHTTP.unixSocket); err != nil {
			remote.Cancel(errors.Wrap(err))
			return
		}
	}()

	return
}

type HTTP struct {
	unixSocket repo_local.UnixSocket
	remote     *repo_local.Repo
}

func (remote *HTTP) GetRepo() repo.Repo {
	return remote
}

func (remote *HTTP) GetBlobStore() dir_layout.BlobStore {
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

func (remote *HTTP) do(
	request *http.Request,
) (response *http.Response, err error) {
	var conn net.Conn

	if conn, err = net.Dial("unix", remote.unixSocket.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	bw := bufio.NewWriter(conn)

	if err = request.Write(bw); err != nil {
		err = errors.Errorf("failed to write to socket: %w", err)
		return
	}

	if err = bw.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if response, err = http.ReadResponse(
		bufio.NewReader(conn),
		request,
	); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}
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

	if response, err = remote.do(request); err != nil {
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

func (blobStore *HTTPBlobStore) GetBlobStore() dir_layout.BlobStore {
	return blobStore
}

func (blobStore *HTTPBlobStore) HasBlob(sh sha.ShaLike) (ok bool) {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			blobStore.remote.remote.Context,
			"HEAD",
			"/blobs",
			strings.NewReader(sh.GetShaLike().GetShaString()),
		); err != nil {
			blobStore.remote.remote.Context.Cancel(errors.Wrap(err))
			return
		}
	}

	var response *http.Response

	{
		var err error

		if response, err = blobStore.remote.do(request); err != nil {
			blobStore.remote.remote.Context.Cancel(errors.Wrap(err))
			return
		}
	}

	ok = response.StatusCode == http.StatusNoContent

	return
}

func (blobStore *HTTPBlobStore) BlobWriter() (w sha.WriteCloser, err error) {
	return
}

func (blobStore *HTTPBlobStore) BlobReader(
	sh sha.ShaLike,
) (r sha.ReadCloser, err error) {
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

	if response, err = blobStore.remote.do(request); err != nil {
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
