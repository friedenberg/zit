package store_browser

import (
	"bufio"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

// TODO remove entirely
type External struct {
	sku.External
}

func (c *External) GetRepoId() ids.RepoId {
	return *(ids.MustRepoId("browser"))
}

func (e *External) GetObjectId() *ids.ObjectId {
	return e.Transacted.GetObjectId()
}

// TODO support updating bookmarks without overwriting. Maybe move to
// toml-bookmark type
func (e *External) SaveBlob(s fs_home.Home) (err error) {
	var aw sha.WriteCloser

	if aw, err = s.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	var item Item

	if err = item.ReadFromExternal(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	tb := sku_fmt.TomlBookmark{
		Url: item.Url.String(),
	}

	func() {
		bw := bufio.NewWriter(aw)
		defer errors.DeferredFlusher(&err, bw)

		enc := toml.NewEncoder(bw)

		if err = enc.Encode(tb); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

	e.Transacted.Metadata.Blob.SetShaLike(aw)

	return
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (a *External) Clone() sku.ExternalLike {
	b := GetExternalPool().Get()
	sku.ExternalResetter.ResetWith(&b.External, &a.External)
	return b
}
