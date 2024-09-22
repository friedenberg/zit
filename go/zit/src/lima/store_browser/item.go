package store_browser

import (
	"slices"
	"strings"

	"code.linenisgreat.com/chrest/go/src/charlie/browser_items"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Item struct {
	browser_items.Item
}

func (i *Item) GetExternalObjectId() sku.ExternalObjectId {
	return &ids.DumbObjectId{
		Value: i.String(),
		Genre: genres.Zettel,
	}
}

func (i *Item) GetGenre() interfaces.Genre {
	return genres.Zettel
}

func (i *Item) String() string {
	return i.GetKey()
}

func (i *Item) GetKey() string {
	return i.Id.String()
}

func (i *Item) GetObjectId() *ids.ObjectId {
	var oid ids.ObjectId
	errors.PanicIfError(oid.SetLeft(i.GetKey()))
	// errors.PanicIfError(oid.SetRepoId("browser"))
	return &oid
}

func (i *Item) GetType() (t ids.Type, err error) {
	if err = t.Set("browser-" + i.Id.Type); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO move below to !toml-bookmark type
func (i Item) GetUrlPathTag() (e ids.Tag, err error) {
	els := strings.Split(i.Url.Hostname(), ".")
	slices.Reverse(els)

	if els[0] == "www" {
		els = els[1:]
	}

	host := strings.Join(els, "-")

	if len(host) == 0 {
		err = errors.Errorf("empty host: %q", els)
		return
	}

	if err = e.Set("%zz-site-" + host); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i Item) GetTai() (t ids.Tai, err error) {
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

func (i Item) GetDescription() (b descriptions.Description, err error) {
	if err = b.Set(i.Title); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Item) WriteToExternal(e *sku.External) (err error) {
	if err = e.ExternalObjectId.SetRaw(i.Id.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Metadata.Type = ids.MustType("!toml-bookmark")

	m := &e.Metadata

	if m.Tai, err = i.GetTai(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e.ExternalType, err = i.GetType(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Metadata.Fields = []object_metadata.Field{
		{
			Value:              i.Id.String(),
			DisableValueQuotes: true,
			ColorType:          string_format_writer.ColorTypeId,
		},
		{
			Key:       "title",
			Value:     i.Title,
			ColorType: string_format_writer.ColorTypeUserData,
		},
		{
			Key:       "url",
			Value:     i.Url.String(),
			ColorType: string_format_writer.ColorTypeUserData,
		},
	}

	// TODO move to !toml-bookmark type
	var t ids.Tag

	if t, err = i.GetUrlPathTag(); err == nil {
		if err = m.AddTagPtr(&t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	err = nil

	return
}

func (i *Item) ReadFromExternal(e *sku.External) (err error) {
	if err = i.Id.Set(e.ExternalObjectId.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, field := range e.Metadata.Fields {
		switch field.Key {
		case "id":
			if field.Value == "" {
				continue
			}

			if err = i.Id.Set(e.ExternalObjectId.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "", "title":
			i.Title = field.Value

		case "url":
			if err = i.Url.Set(field.Value); err != nil {
				err = errors.Wrap(err)
				return
			}

		default:
			err = errors.Errorf("unsupported field type: %s=%q", field.Key, field.Value)
			return
		}
	}

	// err = todo.Implement()
	return
}
