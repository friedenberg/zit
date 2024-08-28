package store_browser

import (
	"bufio"
	"fmt"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
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
	external_state.State
	sku.Transacted
	browser sku.Transacted
	browserItem
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

	var u *url.URL

	if u, err = e.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	tb := sku_fmt.TomlBookmark{
		Url: u.String(),
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

	e.Metadata.Blob.SetShaLike(aw)

	return
}

func (e *External) SetItem(i browserItem, overwrite bool) (err error) {
	e.browserItem = i

	if err = i.WriteToMetadata(&e.browser.Metadata); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Metadata.Tai = e.browser.Metadata.GetTai()

	// if overwrite {
	// 	if err = i.WriteToMetadata(&e.Metadata); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}
	// }

	// TODO make configurable
	e.Metadata.Type = ids.MustType("!toml-bookmark")

	return
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (t *External) GetExternalObjectId() sku.ExternalObjectId {
	return &t.browserItem
}

func (a *External) Clone() sku.ExternalLike {
	b := GetExternalPool().Get()
	sku.TransactedResetter.ResetWith(&b.Transacted, &a.Transacted)
	sku.TransactedResetter.ResetWith(&b.browser, &a.browser)
	b.browserItem = a.browserItem
  b.State = a.State
	return b
}

func (c *External) GetSku() *sku.Transacted {
	return &c.Transacted
}

func (a *External) GetMetadatei() *object_metadata.Metadata {
	return &a.Metadata
}

func (a *External) GetGattung() interfaces.Genre {
	return a.ObjectId.GetGenre()
}

func (a *External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetObjectId(),
		a.GetObjectSha(),
		a.GetBlobSha(),
	)
}

func (a *External) GetBlobSha() interfaces.Sha {
	return &a.Metadata.Blob
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetObjectId())
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
