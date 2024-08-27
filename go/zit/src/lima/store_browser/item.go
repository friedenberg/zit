package store_browser

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"code.linenisgreat.com/chrest/go/chrest/src/charlie/browser_items"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type browserItem struct {
	browser_items.Item
}

func (i *browserItem) getExternalObjectId() sku.ExternalObjectId {
	return i
}

func (i *browserItem) GetExternalObjectId() *ids.ObjectId {
	return i.GetObjectId()
}

func (i *browserItem) GetKey() string {
	return fmt.Sprintf("%s-%s", i.Id.Type, i.Id.Id)
}

func (i *browserItem) GetObjectId() *ids.ObjectId {
	var oid ids.ObjectId
	errors.PanicIfError(oid.SetLeft(i.GetKey()))
	errors.PanicIfError(oid.SetRepoId("browser"))
	return &oid
}

func (i *browserItem) GetType() (t ids.Type, err error) {
	if err = t.Set("browser-" + i.Id.Type); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *browserItem) GetGenre() interfaces.Genre {
	return genres.Zettel
}

func (dst *browserItem) readFromRaw(src map[string]interface{}) (err error) {
	// TODO BrowserId
	dst.Id.Id = src["id"].(string)
	dst.Id.Type = src["type"].(string)
	dst.Url = src["url"].(string)
	dst.Date = src["date"].(string)
	dst.Title, _ = src["title"].(string)
	dst.ExternalId, _ = src["external-id"].(string)
	return
}

func (i browserItem) WriteToMetadata(m *object_metadata.Metadata) (err error) {
	if m.Tai, err = i.GetTai(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.Type, err = i.GetType(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.Description, err = i.GetDescription(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var u *url.URL

	if u, err = i.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var e ids.Tag

	els := strings.Split(u.Hostname(), ".")
	slices.Reverse(els)

	if els[0] == "www" {
		els = els[1:]
	}

	host := strings.Join(els, "-")

	if len(host) > 0 {
		if err = e.Set("zz-site-" + host); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = m.AddTagPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i browserItem) WriteToObjectIds(oids ...*ids.ObjectId) (err error) {
	for _, oid := range oids {
		if err = i.WriteToObjectId(oid); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i browserItem) WriteToObjectId(oi *ids.ObjectId) (err error) {
	oi.ResetWith(i.GetObjectId())
	return
}

func (i browserItem) GetTai() (t ids.Tai, err error) {
	if i.Date == "" {
		return
	}

	if err = t.SetFromRFC3339(i.Date); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

var errEmptyUrl = errors.New("empty url")

func (i browserItem) GetUrl() (u *url.URL, err error) {
	ur := i.Url

	if ur == "" {
		err = errEmptyUrl
		return
	}

	if u, err = url.Parse(ur); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i browserItem) GetDescription() (b descriptions.Description, err error) {
	if err = b.Set(i.Title); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
