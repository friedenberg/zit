package sku_fmt

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type JsonWithUrl struct {
	Json
	TomlBookmark
}

func MakeJsonTomlBookmark(
	sk *sku.Transacted,
	s repo_layout.Layout,
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

	if _, err = url.Parse(j.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, tabRaw := range tabs {
		tab := tabRaw.(map[string]interface{})

		if _, err = url.Parse(tab["url"].(string)); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
