package remote_http

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/india/log_remote_inventory_lists"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

func MakeClient(
	envUI env_ui.Env,
	transport http.RoundTripper,
	localRepo repo.LocalRepo,
	typedBlobStore typed_blob_store.InventoryList,
) *client {
	client := &client{
		envUI: envUI,
		http: http.Client{
			Transport: transport,
		},
		localRepo:      localRepo,
		typedBlobStore: typedBlobStore,
	}

	client.Initialize()

	return client
}

type client struct {
	envUI           env_ui.Env
	configImmutable config_immutable_io.ConfigLoadedPublic
	http            http.Client
	localRepo       repo.LocalRepo
	typedBlobStore  typed_blob_store.InventoryList

	logRemoteInventoryLists log_remote_inventory_lists.Log
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

	decoder := config_immutable_io.CoderPublic{}

	if _, err := decoder.DecodeFrom(
		&client.configImmutable,
		response.Body,
	); err != nil {
		client.envUI.CancelWithErrorAndFormat(err, "failed to read remote immutable config")
	}

	client.logRemoteInventoryLists = log_remote_inventory_lists.Make(
		client.localRepo.GetEnvRepo(),
		client.localRepo.GetEnvRepo(),
	)
}

func (client *client) GetEnv() env_ui.Env {
	return client.envUI
}

func (client *client) GetImmutableConfigPublic() config_immutable_io.ConfigLoadedPublic {
	return client.configImmutable
}

func (client *client) GetInventoryListStore() sku.InventoryListStore {
	return client
}

func (client *client) GetTypedInventoryListBlobStore() typed_blob_store.InventoryList {
	return client.typedBlobStore
}

func (client *client) GetBlobStore() interfaces.BlobStore {
	return client
}

func (client *client) GetObjectStore() sku.ObjectStore {
	return nil
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
	builderOptions query.BuilderOption,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Query, err error) {
	err = todo.Implement()
	return
}

func (remote *client) MakeInventoryList(
	queryGroup *query.Query,
) (list *sku.List, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		remote.GetEnv(),
		"GET",
		// fmt.Sprintf("/query/%s", queryGroup.String()),
		fmt.Sprintf("/query/%s", url.QueryEscape(queryGroup.String())),
		nil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = remote.http.Do(request); err != nil {
		err = errors.ErrorWithStackf("failed to read response: %w", err)
		return
	}

	if err = ReadErrorFromBodyOnNot(response, 200); err != nil {
		err = errors.Wrap(err)
		return
	}

	listFormat := remote.GetInventoryListStore().FormatForVersion(
		remote.GetImmutableConfigPublic().ImmutableConfig.GetStoreVersion(),
	)

	list = sku.MakeList()

	if err = inventory_list_blobs.ReadInventoryListBlob(
		listFormat,
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
	queryGroup *query.Query,
	options repo.RemoteTransferOptions,
) (err error) {
	return client.pullQueryGroupFromWorkingCopy(
		remote.(repo.WorkingCopy),
		queryGroup,
		options,
	)
}

func (client *client) pullQueryGroupFromWorkingCopy(
	remote repo.WorkingCopy,
	queryGroup *query.Query,
	options repo.RemoteTransferOptions,
) (err error) {
	var list *sku.List

	if list, err = remote.MakeInventoryList(queryGroup); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO local / remote version negotiation

	listFormat := client.GetInventoryListStore().FormatForVersion(
		config_immutable.CurrentStoreVersion,
	)

	buffer := bytes.NewBuffer(nil)

	for {
		client.envUI.ContinueOrPanicOnDone()

		// TODO make a reader version of inventory lists to avoid allocation
		if _, err = listFormat.WriteInventoryListBlob(list, buffer); err != nil {
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
				buffer,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if options.AllowMergeConflicts {
				request.Header.Add("x-zit-remote_transfer_options-allow_merge_conflicts", "true")
			}

			if response, err = client.http.Do(request); err != nil {
				err = errors.ErrorWithStackf("failed to read response: %w", err)
				return
			}
		}

		if err = ReadErrorFromBodyOnNot(
			response,
			http.StatusCreated,
			http.StatusExpectationFailed,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		br := bufio.NewReader(response.Body)

		client.GetEnv().ContinueOrPanicOnDone()

		var listMissingSkus *sku.List

		if listMissingSkus, err = client.typedBlobStore.ReadInventoryListBlob(
			builtin_types.GetOrPanic(builtin_types.InventoryListTypeV1).Type,
			br,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = response.Body.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}

		// if options.IncludeBlobs {
		for expected := range listMissingSkus.All() {
			client.envUI.ContinueOrPanicOnDone()

			if err = client.WriteBlobToRemote(
				remote.GetBlobStore(),
				sha.Make(expected.GetBlobSha()),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
		// }

		if response.StatusCode == http.StatusCreated {
			ui.Log().Print("done")
			return
		}

		buffer.Reset()
	}

	return
}

func (client *client) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	err = todo.Implement()
	return
}
