package sku_fmt

import (
	"bytes"
	"io"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type TomlBookmark struct {
	Url string `toml:"url"`
}

func TomlBookmarkUrl(
	sk *sku.Transacted,
	s fs_home.Standort,
) (ur *url.URL, err error) {
	var r sha.ReadCloser

	if r, err = s.BlobReader(sk.GetAkteSha()); err != nil {
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
