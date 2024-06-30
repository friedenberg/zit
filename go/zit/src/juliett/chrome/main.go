package chrome

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"syscall"
	"time"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
)

type transacted struct {
	sync.Mutex
	schnittstellen.MutableSetLike[*kennung.Kennung2]
}

type Store struct {
	konfig       *konfig.Compiled
	storeFuncs   sku.StoreFuncs
	typ          kennung.Typ
	chrestConfig chrest.Config
	standort     standort.Standort

	urls map[url.URL][]item

	l       sync.Mutex
	removed map[url.URL]struct{}
	added   map[url.URL]struct{}

	transacted transacted

	transactedUrlIndex map[url.URL]sku.TransactedMutableSet

	itemDeletedStringFormatWriter schnittstellen.FuncIter[*CheckedOut]
}

func MakeChrome(
	k *konfig.Compiled,
	s standort.Standort,
	itemDeletedStringFormatWriter schnittstellen.FuncIter[*CheckedOut],
) *Store {
	c := &Store{
		konfig:   k,
		typ:      kennung.MustTyp("toml-bookmark"),
		standort: s,
		removed:  make(map[url.URL]struct{}),
		added:    make(map[url.URL]struct{}),
		transacted: transacted{
			MutableSetLike: collections_value.MakeMutableValueSet(
				iter.StringerKeyer[*kennung.Kennung2]{},
			),
		},
		transactedUrlIndex:            make(map[url.URL]sku.TransactedMutableSet),
		itemDeletedStringFormatWriter: itemDeletedStringFormatWriter,
	}

	return c
}

func (c *Store) GetVirtualStore() sku.ExternalStoreLike {
	return c
}

func (c *Store) Flush() (err error) {
	if c.konfig.DryRun {
		return
	}

	if len(c.removed) == 0 && len(c.added) == 0 {
		return
	}

	var req *http.Request

	if req, err = http.NewRequest("PUT", "http://localhost/urls", nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := bytes.NewBuffer(nil)
	var reqPayload requestUrlsPut
	reqPayload.Deleted = make([]string, 0, len(c.removed))

	for u := range c.removed {
		reqPayload.Deleted = append(reqPayload.Deleted, u.String())
	}

	for u := range c.added {
		reqPayload.Added = append(
			reqPayload.Added,
			createOneTabRequest{
				Url: u.String(),
			},
		)
	}

	enc := json.NewEncoder(b)

	if err = enc.Encode(reqPayload); err != nil {
		err = errors.Wrap(err)
		return
	}

	req.Body = io.NopCloser(b)

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Duration(1e9)),
	)

	defer cancel()

	if _, err = chrest.AskChrome(ctx, c.chrestConfig, req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			ui.Err().Print("chrest offline")
			err = nil
		} else if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	clear(c.added)
	clear(c.removed)

	return
}

func (c *Store) getUrl(sk *sku.Transacted) (u *url.URL, err error) {
	var r sha.ReadCloser

	if r, err = c.standort.AkteReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var tb sku_fmt.TomlBookmark

	dec := toml.NewDecoder(r)

	if err = dec.Decode(&tb); err != nil {
		err = errors.Wrapf(err, "Sha: %s, Kennung: %s", sk.GetAkteSha(), sk.GetKennung())
		return
	}

	if u, err = url.Parse(tb.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) QueryCheckedOut(
	qg sku.ExternalQuery,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	// o := sku.ObjekteOptions{
	// 	Mode: objekte_mode.ModeRealizeSansProto,
	// }

	var co CheckedOut

	for u, items := range c.urls {
		existing, ok := c.transactedUrlIndex[u]

		if !ok && qg.ExcludeUntracked {
			continue
		}

		for _, item := range items {
			processOne := func(internal *sku.Transacted) (err error) {
				co.External.browser.Kennung.SetGattung(gattung.Zettel)
				co.External.Kennung.SetGattung(gattung.Zettel)

				if internal != nil {
					sku.TransactedResetter.ResetWith(&co.Internal, internal)
					sku.TransactedResetter.ResetWith(&co.External.Transacted, internal)
					co.State = checked_out_state.StateExistsAndSame
				} else {
					sku.TransactedResetter.Reset(&co.External.Transacted)
					sku.TransactedResetter.Reset(&co.Internal)
					co.State = checked_out_state.StateUntracked
				}

				browser := &co.External.browser
				co.External.item = item

				if co.External.Metadatei.Tai, err = item.GetTai(); err != nil {
					err = errors.Wrap(err)
					return
				}

				browser.Metadatei.Tai = co.External.Metadatei.Tai

				if browser.Metadatei.Typ, err = item.GetTyp(); err != nil {
					err = errors.Wrap(err)
					return
				}

				if browser.Metadatei.Bezeichnung, err = item.GetBezeichnung(); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !qg.ContainsSku(browser) {
					return
				}

				if err = f(&co); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			}

			if !ok {
				if err = processOne(nil); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				if err = existing.Each(processOne); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}
	}

	return
}

func (c *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	if !sz.Metadatei.Typ.Equals(c.typ) {
		err = errors.Wrap(sku.ErrExternalStoreUnsupportedTyp(sz.Metadatei.Typ))
		return
	}

	var u *url.URL

	if u, err = c.getUrl(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	co := GetCheckedOutPool().Get()
	cz = co

	sku.TransactedResetter.ResetWith(co.GetSku(), sz)
	sku.TransactedResetter.ResetWith(co.GetSkuExternalLike().GetSku(), sz)
	co.State = checked_out_state.StateJustCheckedOut
	co.External.browser.Metadatei.Typ = kennung.MustTyp("!chrome-tab")
	co.External.browser.Metadatei.Typ = kennung.MustTyp("!chrome-tab")
	co.External.item = map[string]interface{}{"url": u.String()}

	c.l.Lock()
	defer c.l.Unlock()

	c.added[*u] = struct{}{}

	// 	ui.Debug().Print(response)

	return
}
