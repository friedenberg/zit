package verzeichnisse

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

const digitWidth = 2
const pageCount = 1 << (digitWidth * 4)

type Zettelen struct {
	konfig.Konfig
	path string
	pool ZettelPool
	ioFactory
	pages [pageCount]*zettelenPageWithState
}

func MakeZettelen(
	k konfig.Konfig,
	s standort.Standort,
	f ioFactory,
	p zettel_transacted.Pool,
) (i *Zettelen, err error) {
	i = &Zettelen{
		Konfig:    k,
		path:      s.DirVerzeichnisseZettelenNeue(),
		ioFactory: f,
		pool:      MakeZettelPool(),
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

	z := i.pool.Get()
	z.Transacted = tz
	z.EtikettenExpandedSorted = tz.Named.Stored.Zettel.Etiketten.Expanded().SortedString()
	z.EtikettenSorted = tz.Named.Stored.Zettel.Etiketten.SortedString()

	if err = p.Add(z); err != nil {
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
	w1 Writer,
	qs ...zettel_named.NamedFilter,
) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(i.pages))

	w := writer{
		writers: []Writer{
			MakeWriter(i.shouldIncludeVerzeichnisse),
			WriterZettelTransacted{
				Writer: zettel_transacted.MakeWriter(
					func(zt *zettel_transacted.Zettel) (err error) {
						for _, q := range qs {
							if !q.IncludeNamedZettel(zt.Named) {
								err = io.EOF
								return
							}
						}

						return
					},
				),
			},
			w1,
		},
		ZettelPool: &i.pool,
	}

	for _, p := range i.pages {
		go func(p *zettelenPageWithState) {
			defer wg.Done()

			if err = p.WriteZettelenTo(w); err != nil {
				err = errors.Wrap(err)
				return
			}
		}(p)
	}

	wg.Wait()

	return
}

func (i *Zettelen) shouldIncludeVerzeichnisse(z *Zettel) (err error) {
	if i.IncludeHidden {
		return
	}

	for _, p := range z.EtikettenExpandedSorted {
		for tn, tv := range i.Tags {
			if !tv.Hide {
				continue
			}

			if strings.HasPrefix(p, tn) {
				err = io.EOF
				return
			}
		}
	}

	return
}
