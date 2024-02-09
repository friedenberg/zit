package commands

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

func (c CheckinTomlBookmark) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet()
}

func (c CheckinTomlBookmark) Run(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	rb := bufio.NewReader(u.In())

	urlsFound := make(map[string]struct{})

	for {
		var line string

		line, err = rb.ReadString('\n')

		if errors.IsNotNilAndNotEOF(err) {
			err = errors.Wrap(err)
			return
		}

		isEOF := err == io.EOF

		line = strings.TrimSpace(line)

		if line != "" {
			var ur *url.URL

			if ur, err = url.Parse(line); err != nil {
				err = errors.Wrap(err)
				return
			}

			urlsFound[ur.String()] = struct{}{}
		}

		if isEOF {
			break
		}
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
		if err = etiketten.EachPtr(swu.Metadatei.AddEtikettPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		if u.StoreObjekten().CreateOrUpdate(swu, swu.GetKennung()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for ur := range urlsFound {
		content := fmt.Sprintf("url = \"%s\"", ur)

		mg := metadatei.GetPool().Get()
		mg.SetEtiketten(etiketten)

		if err = mg.Typ.Set("toml-bookmark"); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = u.StoreObjekten().CreateWithAkteString(
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
	urlsFound map[string]struct{},
) (urls map[string]SkuWithUrl, err error) {
	query := "!toml-bookmark:z"

	ids := u.MakeMetaIdSetWithExcludedHidden(c.DefaultGattungen())

	if err = ids.Set(query); err != nil {
		err = errors.Wrap(err)
		return
	}

	urls = make(map[string]SkuWithUrl)

	if err = u.StoreObjekten().QueryWithoutCwd(
		ids,
		iter.MakeSyncSerializer(
			func(sk *sku.Transacted) (err error) {
				var url *url.URL

				if url, err = sku_fmt.TomlBookmarkUrl(sk, u.Standort()); err != nil {
					err = errors.Wrap(err)
					return
				}

				urlString := url.String()

				if _, ok := urlsFound[urlString]; !ok {
					return
				}

				delete(urlsFound, urlString)

				sk2 := sku.GetTransactedPool().Get()
				sku.TransactedResetter.ResetWith(sk2, sk)

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
