package store_browser

import (
	"bufio"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

type External struct {
	sku.External
	Item Item
}

func (c *External) GetRepoId() ids.RepoId {
	return *(ids.MustRepoId("browser"))
}

func (e *External) GetObjectId() *ids.ObjectId {
	return e.Transacted.GetObjectId()
}

func (e *External) GetExternalState() external_state.State {
	return e.State
}

func (e *External) SaveBlob(s fs_home.Home) (err error) {
	var aw sha.WriteCloser

	if aw, err = s.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	tb := sku_fmt.TomlBookmark{
		Url: e.Item.Url.String(),
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

func (e *External) SetItem(i Item, overwrite bool) (err error) {
	e.Item = i

  m := &e.Transacted.Metadata

	if m.Tai, err = i.GetTai(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e.ExternalType, err = i.GetType(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.Description, err = i.GetDescription(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var t ids.Tag

	if t, err = i.GetUrlPathTag(); err == nil {
		if err = m.AddTagPtr(&t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	err = nil

	e.Transacted.Metadata.Type = ids.MustType("!toml-bookmark")

	return
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (t *External) GetExternalObjectId() sku.ExternalObjectId {
	return t.Item.GetExternalObjectId()
}

func (a *External) Clone() sku.ExternalLike {
	b := GetExternalPool().Get()
	sku.TransactedResetter.ResetWith(b.GetSku(), a.GetSku())
	b.Item = a.Item
	b.State = a.State
	return b
}

func (a *External) GetMetadatei() *object_metadata.Metadata {
	return &a.Transacted.Metadata
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGenre(), o.GetObjectId())
}

type lessorExternal struct{}

func (lessorExternal) Less(a, b External) bool {
	panic("not supported")
}

func (lessorExternal) LessPtr(a, b *External) bool {
	return a.Transacted.GetTai().Less(b.Transacted.GetTai())
}

type equalerExternal struct{}

func (equalerExternal) Equals(a, b External) bool {
	panic("not supported")
}
