package sku_fmt

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type JsonWithUrl struct {
	Json
	TomlBookmark
}

func MakeJsonTomlBookmark(
	sk *sku.Transacted,
	s dir_layout.DirLayout,
	tabs []interface{},
) (j JsonWithUrl, err error) {
	if err = j.FromTransacted(sk, s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = toml.Unmarshal([]byte(j.BlobString), &j.TomlBookmark); err != nil {
		err = errors.Wrapf(err, "%q", j.BlobString)
		return
	}

	var u1 *url.URL

	if u1, err = url.Parse(j.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, tabRaw := range tabs {
		tab := tabRaw.(map[string]interface{})
		var u *url.URL

		if u, err = url.Parse(tab["url"].(string)); err != nil {
			err = errors.Wrap(err)
			return
		}

		if *u == *u1 {
			ui.Debug().Print(u, u1)
		}
	}

	return
}
