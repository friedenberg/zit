package sku_fmt

import (
	"bytes"
	"io"
	"net/url"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type TomlBookmark struct {
	Url string `toml:"url"`
}

func TomlBookmarkUrl(
	sk *sku.Transacted,
	s standort.Standort,
) (ur *url.URL, err error) {
	var r sha.ReadCloser

	if r, err = s.AkteReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var out bytes.Buffer

	if _, err = io.Copy(&out, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	var tb TomlBookmark

	if err = toml.Unmarshal(out.Bytes(), &tb); err != nil {
		err = errors.Wrapf(err, "%q", out.String())
		return
	}

	if ur, err = url.Parse(tb.Url); err != nil {
		err = errors.Wrapf(err, "%q", tb.Url)
		return
	}

	return
}
