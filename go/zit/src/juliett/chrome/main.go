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

func (c *Chrome) GetVirtualStore() sku.VirtualStore {
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

	if req, err = http.NewRequest("DELETE", "http://localhost/urls", nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := bytes.NewBuffer(nil)
	urls := make([]string, 0, len(c.removed))

	for u := range c.removed {
		urls = append(urls, u.String())
	}

	enc := json.NewEncoder(b)

	if err = enc.Encode(urls); err != nil {
		err = errors.Wrap(err)
		return
	}

	req.Body = io.NopCloser(b)

	if _, err = chrest.AskChrome(c.chrestConfig, req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			errors.Err().Print("chrest offline")
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

func (c *Chrome) CommitTransacted(sk *sku.Transacted) (err error) {
	if !c.transacted.Contains(&sk.Kennung) {
		return
	}

	ees := sk.Metadatei.Verzeichnisse.GetExpandedEtiketten()
	es := sk.Metadatei.GetEtiketten()
	log.Debug().Print(iter.StringCommaSeparated(es), iter.StringCommaSeparated(ees))

	var u *url.URL

	if u, err = c.getUrl(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, ok := c.urls[*u]; !ok {
		// TODO fetch previous URL
		return
	}

	c.removed[*u] = struct{}{}

	return
}

func (c *Chrome) ContainsSku(sk *sku.Transacted) bool {
	if !sk.GetTyp().Equals(c.typ) {
		return false
	}

	u, err := c.getUrl(sk)
	errors.PanicIfError(err)

	ts, ok := c.urls[*u]

	if !ok {
		return false
	}

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

	c.transacted.Lock()
	defer c.transacted.Unlock()

	errors.PanicIfError(c.transacted.Add(&sk.Kennung))

	return true
}
