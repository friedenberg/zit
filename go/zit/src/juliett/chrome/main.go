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
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
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

	urls       map[url.URL][]item
	removed    map[url.URL]struct{}
	transacted transacted

	transactedUrlIndex map[url.URL]sku.TransactedMutableSet
}

func MakeChrome(
	k *konfig.Compiled,
	s standort.Standort,
	storeFuncs sku.StoreFuncs,
) *Store {
	c := &Store{
		konfig:     k,
		storeFuncs: storeFuncs,
		typ:        kennung.MustTyp("toml-bookmark"),
		standort:   s,
		removed:    make(map[url.URL]struct{}),
		transacted: transacted{
			MutableSetLike: collections_value.MakeMutableValueSet(
				iter.StringerKeyer[*kennung.Kennung2]{},
			),
		},
		transactedUrlIndex: make(map[url.URL]sku.TransactedMutableSet),
	}

	return c
}

func (c *Store) GetVirtualStore() query.VirtualStore {
	return c
}

func (c *Store) Flush() (err error) {
	if c.konfig.DryRun || !c.konfig.ChrestEnabled {
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
		}

		return
	}

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
		for _, item := range items {
			processOne := func(internal *sku.Transacted) (err error) {
				if internal != nil {
					sku.TransactedResetter.ResetWith(&co.Internal, internal)
					sku.TransactedResetter.ResetWith(&co.External.Transacted, internal)
					co.State = checked_out_state.StateExistsAndSame
				} else {
					sku.TransactedResetter.Reset(&co.External.Transacted)
					sku.TransactedResetter.Reset(&co.Internal)
					co.State = checked_out_state.StateUntracked
					co.External.Kennung.SetGattung(gattung.Zettel)
					co.External.browser.Kennung.SetGattung(gattung.Zettel)
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

			existing, ok := c.transactedUrlIndex[u]

			if !ok {
				if qg.ExcludeUntracked {
					continue
				}

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
