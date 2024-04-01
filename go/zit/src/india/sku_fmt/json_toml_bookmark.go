package sku_fmt

import (
	"net/url"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type JsonWithUrl struct {
	Json
	TomlBookmark
}

func MakeJsonTomlBookmark(
	sk *sku.Transacted,
	s standort.Standort,
	chromeTabs []interface{},
) (j JsonWithUrl, err error) {
	if err = j.FromTransacted(sk, s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = toml.Unmarshal([]byte(j.Akte), &j.TomlBookmark); err != nil {
		err = errors.Wrapf(err, "%q", j.Akte)
		return
	}

	var u1 *url.URL

	if u1, err = url.Parse(j.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, tabRaw := range chromeTabs {
		tab := tabRaw.(map[string]interface{})
		var u *url.URL

		if u, err = url.Parse(tab["url"].(string)); err != nil {
			err = errors.Wrap(err)
			return
		}

		if *u == *u1 {
			log.Debug().Print(u, u1)
		}
	}

	return
}
