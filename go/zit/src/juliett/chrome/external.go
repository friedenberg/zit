package chrome

import (
	"bufio"
	"fmt"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

type External struct {
	sku.Transacted
	browser sku.Transacted
	item
}

func (e *External) SaveAkte(s standort.Standort) (err error) {
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

	e.Metadatei.Akte.SetShaLike(aw)

	return
}

func (e *External) SetItem(i item, overwrite bool) (err error) {
	e.item = i

	if err = i.WriteToMetadatei(&e.browser.Metadatei); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Metadatei.Tai = e.browser.Metadatei.GetTai()

	if overwrite {
		if err = i.WriteToMetadatei(&e.Metadatei); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO make configurable
	e.Metadatei.Typ = kennung.MustType("!toml-bookmark")

	return
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
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

func (a *External) GetKennung() kennung.IdLike {
	return &a.Kennung
}

func (a *External) GetMetadatei() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a *External) GetGattung() interfaces.Genre {
	return a.Kennung.GetGenre()
}

func (a *External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetKennung(),
		a.GetObjekteSha(),
		a.GetAkteSha(),
	)
}

func (a *External) GetAkteSha() interfaces.ShaLike {
	return &a.Metadatei.Akte
}

// func (a *External) SetAkteSha(v schnittstellen.ShaLike) (err error) {
// 	if err = a.Metadatei.Akte.SetShaLike(v); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = a.FDs.Akte.SetShaLike(v); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
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
