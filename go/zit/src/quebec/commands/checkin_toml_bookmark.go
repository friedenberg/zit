package commands

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CheckinTomlBookmark struct{}

func init() {
	registerCommand(
		"checkin-toml-bookmark",
		func(f *flag.FlagSet) Command {
			c := &CheckinTomlBookmark{}

			return c
		},
	)
}

func (c CheckinTomlBookmark) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung()
}

type CheckinTomlBookmarkEntry struct {
	UrlString string   `json:"url"`
	Url       *url.URL `json:"-"`
	Title     string   `json:"title"`
}

func (c CheckinTomlBookmark) Run(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	rb := bufio.NewReader(u.In())

	urlsFound := make(map[string]CheckinTomlBookmarkEntry)

	dec := json.NewDecoder(rb)

	for {
		var entry CheckinTomlBookmarkEntry

		err = dec.Decode(&entry)

		if errors.IsNotNilAndNotEOF(err) {
			err = errors.Wrap(err)
			return
		}

		isEOF := err == io.EOF

		if isEOF {
			break
		}

		if entry.Url, err = url.Parse(entry.UrlString); err != nil {
			err = errors.Wrap(err)
			return
		}

		urlsFound[entry.UrlString] = entry
	}

	var urls map[string]SkuWithUrl

	if urls, err = c.getUrls(u, urlsFound); err != nil {
		err = errors.Wrap(err)
		return
	}

	var etiketten kennung.EtikettSet

	if etiketten, err = kennung.MakeEtikettSetStrings(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	for _, swu := range urls {
		if err = etiketten.EachPtr(swu.AddEtikettPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		if u.GetStore().CreateOrUpdate(swu, swu.GetKennung()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, entry := range urlsFound {
		content := fmt.Sprintf("url = \"%s\"", entry.UrlString)

		mg := metadatei.GetPool().Get()
		mg.SetEtiketten(etiketten)

		if err = mg.Bezeichnung.Set(entry.Title); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = mg.Typ.Set("toml-bookmark"); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = u.GetStore().CreateWithAkteString(
			mg,
			content,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

type SkuWithUrl struct {
	*sku.Transacted
	*url.URL
}

func (c CheckinTomlBookmark) getUrls(
	u *umwelt.Umwelt,
	urlsFound map[string]CheckinTomlBookmarkEntry,
) (urls map[string]SkuWithUrl, err error) {
	q := "!toml-bookmark?z"

	builder := u.MakeQueryBuilderExcludingHidden(c.DefaultGattungen())
	var ids *query.Group

	if ids, err = builder.BuildQueryGroup(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	urls = make(map[string]SkuWithUrl)

	if err = u.GetStore().QueryWithoutCwd(
		ids,
		iter.MakeSyncSerializer(
			func(sk *sku.Transacted) (err error) {
				var url *url.URL

				if url, err = sku_fmt.TomlBookmarkUrl(sk, u.Standort()); err != nil {
					err = errors.Wrap(err)
					return
				}

				urlString := url.String()

				entry, ok := urlsFound[urlString]

				if !ok {
					return
				}

				title := entry.Title

				delete(urlsFound, urlString)

				sk2 := sku.GetTransactedPool().Get()
				sku.TransactedResetter.ResetWith(sk2, sk)

				if err = sk2.Metadatei.Bezeichnung.Set(title); err != nil {
					err = errors.Wrap(err)
					return
				}

				urls[urlString] = SkuWithUrl{
					Transacted: sk2,
					URL:        url,
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
