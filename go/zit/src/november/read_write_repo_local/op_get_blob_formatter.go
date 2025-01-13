package read_write_repo_local

import (
	"maps"
	"slices"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/india/sku"
)

func (u *Repo) GetBlobFormatter(
	tipe ids.Type,
	formatId string,
	utiGroup string,
) (blobFormatter script_config.RemoteScript, err error) {
	if tipe.GetType().IsEmpty() {
		ui.Err().Print("empty type")
		return
	}

	var typeObject *sku.Transacted

	if typeObject, err = u.GetStore().ReadTransactedFromObjectId(
		tipe.GetType(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var typeBlob type_blobs.Blob

	if typeBlob, _, err = u.GetStore().GetBlobStore().GetType().ParseTypedBlob(
		typeObject.GetType(),
		typeObject.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer u.GetStore().GetBlobStore().GetType().PutTypedBlob(
		typeObject.GetType(),
		typeBlob,
	)

	ok := false

	if utiGroup == "" {
		getBlobFormatter := func(formatId string) script_config.RemoteScript {
			var formatIds []string

			if formatId == "" {
				formatIds = []string{"text-edit", "text"}
			} else {
				formatIds = []string{formatId}
			}

			for _, formatId := range formatIds {
				blobFormatter, ok = typeBlob.GetFormatters()[formatId]

				if ok {
					return blobFormatter
				}
			}

			return nil
		}

		blobFormatter = getBlobFormatter(formatId)

		return
	}

	var g type_blobs.UTIGroup
	g, ok = typeBlob.GetFormatterUTIGroups()[utiGroup]

	if !ok {
		err = errors.BadRequestf(
			"no uti group: %q. Available groups: %s",
			utiGroup,
			maps.Keys(typeBlob.GetFormatterUTIGroups()),
		)
		return
	}

	ft, ok := g.Map()[formatId]

	if !ok {
		err = errors.Errorf(
			"no format id %q for uti group %q. Available groups: %s",
			formatId,
			utiGroup,
			slices.Collect(maps.Keys(g.Map())),
		)

		return
	}

	formatId = ft

	blobFormatter, ok = typeBlob.GetFormatters()[formatId]

	if !ok {
		ui.Err().Print("no matching format id")
		blobFormatter = nil
		// TODO-P2 allow option to error on missing format
		// err = errors.Normalf("no format id %q", actualFormatId)
		// return

		return
	}

	return
}
