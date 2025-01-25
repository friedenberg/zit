package remote_http

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
)

func MakeClient(
	envUI env_ui.Env,
	transport http.RoundTripper,
	localInventoryListStore sku.InventoryListStore,
) *client {
	client := &client{
		envUI: envUI,
		http: http.Client{
			Transport: transport,
		},
		localInventoryListStore: localInventoryListStore,
	}

	client.Initialize()

	return client
}

type client struct {
	envUI                   env_ui.Env
	configImmutable         config_immutable_io.ConfigLoaded
	http                    http.Client
	localInventoryListStore sku.InventoryListStore
}

func (client *client) Initialize() {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			client.GetEnv(),
			"GET",
			"/config-immutable",
			nil,
		); err != nil {
			client.envUI.CancelWithError(err)
		}
	}

	var response *http.Response

	{
		var err error

		if response, err = client.http.Do(request); err != nil {
			client.envUI.CancelWithErrorAndFormat(err, "failed to read response")
		}
	}

	decoder := config_immutable_io.Coder{}

	if _, err := decoder.DecodeFrom(
		&client.configImmutable,
		response.Body,
	); err != nil {
		client.envUI.CancelWithErrorAndFormat(err, "failed to read remote immutable config")
	}
}

func (client *client) GetEnv() env_ui.Env {
	return client.envUI
}

func (client *client) GetImmutableConfig() config_immutable_io.ConfigLoaded {
	return client.configImmutable
}

func (client *client) GetInventoryListStore() sku.InventoryListStore {
	return client
}

func (client *client) GetBlobStore() interfaces.BlobStore {
	return client
}

func (client *client) MakeImporter(
	options sku.ImporterOptions,
	storeOptions sku.StoreOptions,
) sku.Importer {
	panic(todo.Implement())
}

func (client *client) ImportList(
	list *sku.List,
	i sku.Importer,
) (err error) {
	return todo.Implement()
}

func (client *client) MakeExternalQueryGroup(
	builderOptions query.BuilderOptions,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Group, err error) {
	err = todo.Implement()
	return
}

func (client *client) MakeInventoryList(
	qg *query.Group,
) (list *sku.List, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		client.GetEnv(),
		"GET",
		"/inventory_lists",
		strings.NewReader(qg.String()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = client.http.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	bf := client.GetInventoryListStore().FormatForVersion(
		client.GetImmutableConfig().ImmutableConfig.GetStoreVersion(),
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

// func (remoteHTTP *HTTP) PullQueryGroupFromRemote2(
// 	remote repo.ReadWrite,
// 	options repo.RemoteTransferOptions,
// 	queryStrings ...string,
// ) (err error) {
// 	var qg *query.Group

// 	if qg, err = remoteHTTP.MakeQueryGroup(queryStrings...); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = remoteHTTP.PullQueryGroupFromRemote(
// 		remote,
// 		qg,
// 		options,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

func (client *client) PullQueryGroupFromRemote(
	remote repo.Repo,
	qg *query.Group,
	options repo.RemoteTransferOptions,
) (err error) {
	return client.pullQueryGroupFromWorkingCopy(
		remote.(repo.WorkingCopy),
		qg,
		options,
	)
}

func (client *client) pullQueryGroupFromWorkingCopy(
	remote repo.WorkingCopy,
	queryGroup *query.Group,
	options repo.RemoteTransferOptions,
) (err error) {
	var list *sku.List

	if list, err = remote.MakeInventoryList(queryGroup); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO local / remote version negotiation

	bf := client.GetInventoryListStore().FormatForVersion(
		config_immutable.CurrentStoreVersion,
	)

	b := bytes.NewBuffer(nil)

	// TODO make a reader version of inventory lists to avoid allocation
	if _, err = bf.WriteInventoryListBlob(list, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	{
		var request *http.Request

		if request, err = http.NewRequestWithContext(
			client.GetEnv(),
			"POST",
			"/inventory_lists",
			b,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if options.AllowMergeConflicts {
			request.Header.Add("x-zit-remote_transfer_options-allow_merge_conflicts", "true")
		}

		if response, err = client.http.Do(request); err != nil {
			err = errors.Errorf("failed to read response: %w", err)
			return
		}
	}

	if response.StatusCode >= 300 {
		var sb strings.Builder

		if _, err = io.Copy(&sb, response.Body); err != nil {
		}

		err = errors.BadRequestf("remote responded with error: %q", &sb)
		return
	}

	br := bufio.NewReader(response.Body)

	client.GetEnv().ContinueOrPanicOnDone()

	var shas sha.Slice

	if _, err = shas.ReadFrom(br); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = response.Body.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if options.IncludeBlobs {
		for _, expected := range shas {
			if err = client.WriteBlobToRemote(remote, expected); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	ui.Log().Print("done")

	return
}

func (client *client) WriteBlobToRemote(
	local repo.WorkingCopy,
	expected *sha.Sha,
) (err error) {
	var actual sha.Sha

	// Closed by the http client's transport (our roundtripper calling
	// request.Write)
	var rc interfaces.ShaReadCloser

	if rc, err = local.GetBlobStore().BlobReader(
		expected,
	); err != nil {
		if env_dir.IsErrBlobMissing(err) {
			// TODO make an option to collect this error at the present it, and an
			// option to fetch it from another remote store
			ui.Err().Printf("Blob missing locally: %q", expected)
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var request *http.Request

	if request, err = http.NewRequestWithContext(
		client.GetEnv(),
		"POST",
		"/blobs",
		rc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	request.TransferEncoding = []string{"chunked"}

	var response *http.Response

	if response, err = client.http.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	var shString strings.Builder

	if _, err = io.Copy(&shString, response.Body); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = response.Body.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = actual.Set(strings.TrimSpace(shString.String())); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = expected.AssertEqualsShaLike(&actual); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (client *client) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	err = todo.Implement()
	return
}
