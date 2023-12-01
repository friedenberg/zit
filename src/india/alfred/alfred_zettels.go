package alfred

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/srx/bravo/expansion"
)

func (w *Writer) zettelToItem(
	z *sku.Transacted,
	ha func(*kennung.Hinweis) (string, error),
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = z.Metadatei.Bezeichnung.String()

	es := iter.StringCommaSeparated[kennung.Etikett](
		z.Metadatei.GetEtiketten(),
	)

	k := &z.Kennung
	ks := k.StringFromPtr()

	if a.Title == "" {
		a.Title = ks
		a.Subtitle = es
	} else {
		a.Subtitle = fmt.Sprintf("%s: %s", ks, es)
	}

	a.Arg = ks

	mb := alfred.GetPoolMatchBuilder().Get()
	defer alfred.GetPoolMatchBuilder().Put(mb)
	parts := k.PartsStrings()

	mb.AddMatches(ks)
	mb.AddMatchBytes(parts[0].Bytes())
	mb.AddMatchBytes(parts[2].Bytes())
	mb.AddMatches(z.GetMetadatei().Bezeichnung.String())
	mb.AddMatches(z.GetTyp().String())
	z.Metadatei.GetEtiketten().Each(
		func(e kennung.Etikett) (err error) {
			expansion.ExpanderAll.Expand(
				func(v string) (err error) {
					mb.AddMatches(v)
					return
				},
				e.String(),
			)

			return
		},
	)

	t := z.GetTyp()

	expansion.ExpanderAll.Expand(
		func(v string) (err error) {
			mb.AddMatches(v)
			return
		},
		t.String(),
	)

	// if ha != nil {
	// 	var h hinweis.Hinweis
	// 	var err error

	// 	if h, err = ha.AbbreviateHinweis(z.Hinweis); err != nil {
	// 		return ErrorToItem(err)
	// 	}

	// 	mb.AddMatches(h.String())
	// 	mb.AddMatches(h.Kopf())
	// 	mb.AddMatches(h.Schwanz())
	// }

	a.Match.Write(mb.Bytes())

	// if len(a.Match) > 100 {
	// 	a.Match = a.Match[:100]
	// }

	a.Text.Copy = ks
	a.Uid = "zit://" + ks

	return
}

func (w *Writer) etikettToItem(
	e *kennung.Etikett,
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = "@" + e.String()

	a.Arg = e.String()

	mb := alfred.GetPoolMatchBuilder().Get()
	defer alfred.GetPoolMatchBuilder().Put(mb)

	mb.AddMatches(a.Title)
	mb.AddMatches(iter.Strings[kennung.Etikett](kennung.ExpandOne(e))...)

	a.Match.ReadFromBuffer(&mb.Buffer)

	a.Text.Copy = e.String()
	a.Uid = "zit://" + e.String()

	return
}

func (w *Writer) errorToItem(err error) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = err.Error()

	return
}

func (w *Writer) hinweisToItem(e kennung.Hinweis) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = e.String()
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(),
	// strings.Join(EtikettenStringsFromZettel(z, false), ", "))

	a.Arg = e.String()

	mb := alfred.GetPoolMatchBuilder().Get()
	defer alfred.GetPoolMatchBuilder().Put(mb)

	mb.AddMatch(e.String())
	mb.AddMatch(e.Kopf())
	mb.AddMatch(e.Schwanz())

	a.Match.ReadFromBuffer(&mb.Buffer)

	a.Text.Copy = e.String()
	a.Uid = "zit://" + e.String()

	return
}
