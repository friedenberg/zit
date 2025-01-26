package remote_http

import (
	"iter"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (client client) FormatForVersion(
	sv interfaces.StoreVersion,
) sku.ListFormat {
	return client.localInventoryListStore.FormatForVersion(sv)
}

func (client client) WriteInventoryListObject(t *sku.Transacted) (err error) {
	return todo.Implement()
}

func (client client) ImportInventoryList(
	bs interfaces.BlobStore,
	t *sku.Transacted,
) (err error) {
	return todo.Implement()
}

func (client client) ReadLast() (max *sku.Transacted, err error) {
	return nil, todo.Implement()
}

func (client client) StreamInventoryList(
	blobSha interfaces.Sha,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return todo.Implement()
}

func (client client) ReadAllSkus(
	f func(besty, sk *sku.Transacted) error,
) (err error) {
	return todo.Implement()
}

func (client client) AllInventoryLists() iter.Seq[quiter.ElementOrError[*sku.Transacted]] {
	// var request *http.Request

	// {
	// 	var err error

	// 	if request, err = http.NewRequestWithContext(
	// 		client.GetEnv(),
	// 		"GET",
	// 		"/inventory_lists",
	// 		nil,
	// 	); err != nil {
	// 		client.envUI.CancelWithError(err)
	// 		return nil
	// 	}
	// }

	// var response *http.Response

	// {
	// 	var err error

	// 	if response, err = client.http.Do(request); err != nil {
	// 		client.envUI.CancelWithErrorAndFormat(err, "failed to read response")
	// 		return nil
	// 	}
	// }

	// if err = client.typedBlobStore.DecodeObjectStreamFrom(
	// 	output,
	// 	response.Body,
	// ); err != nil {
	// 	client.envUI.CancelWithError(err)
	// 	return
	// }

	return nil
}

func (client client) ReadAllInventoryLists(
	output interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			client.GetEnv(),
			"GET",
			"/inventory_lists",
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

	if err = client.typedBlobStore.DecodeObjectStreamFrom(
		output,
		response.Body,
	); err != nil {
		client.envUI.CancelWithError(err)
		return
	}

	return
}
