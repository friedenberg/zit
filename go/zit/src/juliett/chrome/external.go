package chrome

import (
	"bufio"
	"fmt"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

type External struct {
	sku.Transacted
	browser sku.Transacted
	item
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

func (e *External) SetItem(i item, overwrite bool) (err error) {
	e.item = i

	if err = i.WriteToMetadata(&e.browser.Metadata); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Metadata.Tai = e.browser.Metadata.GetTai()

	if overwrite {
		if err = i.WriteToMetadata(&e.Metadata); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO make configurable
	e.Metadata.Type = ids.MustType("!toml-bookmark")

	return
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (a *External) Clone() sku.ExternalLike {
  b := GetExternalPool().Get()
  sku.TransactedResetter.ResetWith(&b.Transacted, &a.Transacted)
  sku.TransactedResetter.ResetWith(&b.browser, &a.browser)
  b.item = a.item
	return b
}

func (c *External) GetSku() *sku.Transacted {
	return &c.Transacted
}

func (t *External) SetFromSkuLike(sk sku.SkuLike) (err error) {
	// switch skt := sk.(type) {
	// case *External:
	// TODO reset item with other item
	// }

	if err = t.Transacted.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) GetObjectId() ids.IdLike {
	return &a.ObjectId
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

func (equalerExternal) EqualsPtr(a, b *External) bool {
	return a.EqualsSkuLikePtr(&b.Transacted)
}
