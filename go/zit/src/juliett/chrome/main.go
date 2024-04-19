package chrome

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"syscall"

	"code.linenisgreat.com/chrest"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/juliett/konfig"
)

type transacted struct {
	sync.Mutex
	schnittstellen.MutableSetLike[*kennung.Kennung2]
}

type Chrome struct {
	konfig       *konfig.Compiled
	typ          kennung.Typ
	chrestConfig chrest.Config
	standort     standort.Standort
	urls         map[url.URL][]item
	removed      map[url.URL]struct{}
	transacted   transacted
}

func MakeChrome(k *konfig.Compiled, s standort.Standort) *Chrome {
	c := &Chrome{
		konfig:   k,
		typ:      kennung.MustTyp("toml-bookmark"),
		standort: s,
		removed:  make(map[url.URL]struct{}),
		transacted: transacted{
			MutableSetLike: collections_value.MakeMutableValueSet(
				iter.StringerKeyer[*kennung.Kennung2]{},
			),
		},
	}

	return c
}

func (c *Chrome) GetVirtualStore() query.VirtualStore {
	return c
}

func (c *Chrome) Initialize() (err error) {
	if err = c.chrestConfig.Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var chromeTabsRaw interface{}
	var req *http.Request

	if req, err = http.NewRequest("GET", "http://localhost/urls", nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	if chromeTabsRaw, err = chrest.AskChrome(c.chrestConfig, req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			errors.Err().Print("chrest offline")
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	chromeTabsRaw2 := chromeTabsRaw.([]interface{})

	chromeTabs := make(map[url.URL][]item, len(chromeTabsRaw2))

	for _, tabRaw := range chromeTabsRaw2 {
		tab := tabRaw.(map[string]interface{})
		ur := tab["url"]

		if ur == nil {
			continue
		}

		var u *url.URL

		if u, err = url.Parse(ur.(string)); err != nil {
			err = errors.Wrap(err)
			return
		}

		chromeTabs[*u] = append(chromeTabs[*u], tab)
	}

	c.urls = chromeTabs

	return
}

func (c *Chrome) Flush() (err error) {
	if c.konfig.DryRun {
		return
	}

	if len(c.removed) == 0 {
		return
	}

	var req *http.Request

	if req, err = http.NewRequest("PUT", "http://localhost/urls", nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := bytes.NewBuffer(nil)
	var reqPayload putRequest
	reqPayload.Deleted = make([]string, 0, len(c.removed))

	for u := range c.removed {
		reqPayload.Deleted = append(reqPayload.Deleted, u.String())
	}

	enc := json.NewEncoder(b)

	if err = enc.Encode(reqPayload); err != nil {
		err = errors.Wrap(err)
		return
	}

	req.Body = io.NopCloser(b)

	if _, err = chrest.AskChrome(c.chrestConfig, req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			errors.Err().Print("chrest offline")
			err = nil
		} else if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (c *Chrome) getUrl(sk *sku.Transacted) (u *url.URL, err error) {
	var r sha.ReadCloser

	if r, err = c.standort.AkteReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var tb sku_fmt.TomlBookmark

	dec := toml.NewDecoder(r)

	if err = dec.Decode(&tb); err != nil {
		errors.Err().Print(errors.Wrapf(err, "Sha: %s, Kennung: %s", sk.GetAkteSha(), sk.GetKennung()))
		err = nil
		return
	}

	if u, err = url.Parse(tb.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Chrome) CommitTransacted(kinder, mutter *sku.Transacted) (err error) {
	var dt diff

	if dt, err = c.getDiff(kinder, mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if dt.diffType == diffTypeIgnore {
		return
	}

	var u *url.URL

	if u, err = c.getUrl(kinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, ok := c.urls[*u]; !ok {
		// TODO fetch previous URL
		return
	}

	switch dt.diffType {
	case diffTypeDelete:
		log.Debug().Print("deleted", "TODO add to dedicated printer", kinder)
		c.removed[*u] = struct{}{}

	default:
		log.Debug().Print("TODO not implemented", dt, kinder, mutter)
	}

	return
}

func (c *Chrome) modifySku(sk *sku.Transacted) (didModify bool, err error) {
	if !sk.GetTyp().Equals(c.typ) {
		return
	}

	u, err := c.getUrl(sk)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	ts, ok := c.urls[*u]

	if !ok {
		return
	}

	didModify = true

	for _, t := range ts {
		es := t.Etiketten()

		if err = t.Etiketten().EachPtr(sk.Metadatei.AddEtikettPtr); err != nil {
			errors.PanicIfError(err)
		}

		ex := kennung.ExpandMany(es, expansion.ExpanderRight)

		if err = ex.EachPtr(sk.Metadatei.Verzeichnisse.AddEtikettExpandedPtr); err != nil {
			errors.PanicIfError(err)
		}
	}

	return
}

func (c *Chrome) Query(
	ms *query.Group,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	var sk sku.Transacted

	for _, items := range c.urls {
		for _, item := range items {
			sku.TransactedResetter.Reset(&sk)

			if err = item.HydrateSku(&sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !ms.ContainsSku(&sk) {
				continue
			}

			if err = f(&sk); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (c *Chrome) ModifySku(sk *sku.Transacted) (err error) {
	if _, err = c.modifySku(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Chrome) ContainsSku(sk *sku.Transacted) bool {
	ok, err := c.modifySku(sk)
	if err != nil {
		log.Err().Print(err)
	}

	return ok
}
