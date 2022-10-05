package verzeichnisse

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

const digitWidth = 1
const pageCount = 16 ^ digitWidth

type Zettelen struct {
	konfig.Konfig
	path string
	ioFactory
	pages [pageCount]*zettelenPageWithState
}

func MakeZettelen(
	k konfig.Konfig,
	s standort.Standort,
	f ioFactory,
) (i *Zettelen, err error) {
	i = &Zettelen{
		Konfig:    k,
		path:      s.DirVerzeichnisseZettelenNeue(),
		ioFactory: f,
	}

	for n, _ := range i.pages {
		i.pages[n] = makeZettelenPage(
			i,
			filepath.Join(i.path, fmt.Sprintf("%x", n)),
		)
	}

	return
}

func (i Zettelen) ValidatePageIndex(n int) (err error) {
	switch {
	case n > pageCount:
		fallthrough

	case n < 0:
		err = errors.Errorf("expected page between 0 and %d, but got %d", pageCount-1, n)
		return
	}

	return
}

func (i Zettelen) PageForHinweis(h hinweis.Hinweis) (n int) {
	s := h.Sha()
	ss := s.String()[:digitWidth]

	var err error
	var n1 int64

	if n1, err = strconv.ParseInt(ss, 16, 64); err != nil {
		n1 = -1
	}

	n = int(n1)

	return
}

func (i Zettelen) PagesForZettelTransacted(zt zettel_transacted.Zettel) (ns []int) {
	ns = append(
		ns,
		i.PageForHinweis(zt.Named.Hinweis),
	)

	return
}

func (i *Zettelen) Flush() (err error) {
	errors.Print("flushing")

	for _, p := range i.pages {
		if err = p.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *Zettelen) Add(tz zettel_transacted.Zettel) (err error) {
	ns := i.PagesForZettelTransacted(tz)

	if len(ns) < 1 {
		err = errors.Errorf("expected at least one page value, but got none")
		return
	}

	n := ns[0]

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	if err = p.Add(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) Read(h hinweis.Hinweis) (tz zettel_transacted.Zettel, err error) {
	n := i.PageForHinweis(h)

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	if tz, err = p.ReadHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) ReadMany(
	w zettel_transacted.Writer,
	qs ...zettel_named.NamedFilter,
) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(i.pages))

	for _, p := range i.pages {
		go func(p *zettelenPageWithState) {
			defer wg.Done()

			if err = p.ReadAll(); err != nil {
				err = errors.Wrap(err)
				return
			}

			err = p.innerSet.Each(
				func(tz zettel_transacted.Zettel) (err error) {
					for _, q := range qs {
						if !q.IncludeNamedZettel(tz.Named) {
							return
						}
					}

					if !i.shouldIncludeTransacted(tz) {
						return
					}

					w.WriteZettel(tz)

					return
				},
			)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}(p)
	}

	wg.Wait()

	return
}

func (i *Zettelen) shouldIncludeTransacted(tz zettel_transacted.Zettel) bool {
	if i.IncludeHidden {
		return true
	}

	prefixes := tz.Named.Stored.Zettel.Etiketten.Expanded(etikett.ExpanderRight{})

	for tn, tv := range i.Tags {
		if !tv.Hide {
			continue
		}

		if prefixes.ContainsString(tn) {
			return false
		}
	}

	return true
}
