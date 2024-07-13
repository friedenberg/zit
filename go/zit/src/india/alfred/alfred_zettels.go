package alfred

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/echo/alfred"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (w *Writer) addCommonMatches(
	z *sku.Transacted,
	a *alfred.Item,
) {
	k := &z.Kennung
	ks := k.StringFromPtr()

	mb := alfred.GetPoolMatchBuilder().Get()
	defer alfred.GetPoolMatchBuilder().Put(mb)

	parts := k.PartsStrings()

	mb.AddMatches(ks)
	mb.AddMatchBytes(parts.Left.Bytes())
	mb.AddMatchBytes(parts.Right.Bytes())

	errors.PanicIfError(w.abbr.AbbreviateHinweisOnly(k))
	mb.AddMatches(k.StringFromPtr())
	parts = k.PartsStrings()
	mb.AddMatchBytes(parts.Left.Bytes())
	mb.AddMatchBytes(parts.Right.Bytes())

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

	a.Match.Write(mb.Bytes())
	// a.Match.ReadFromBuffer(&mb.Buffer)
}

func (w *Writer) zettelToItem(
	z *sku.Transacted,
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = z.Metadatei.Bezeichnung.String()

	es := iter.StringCommaSeparated(
		z.Metadatei.GetEtiketten(),
	)

	k := &z.Kennung
	ks := k.StringFromPtr()

	if a.Title == "" {
		a.Title = ks
		a.Subtitle = es
	} else {
		a.Subtitle = fmt.Sprintf("%s: %s %s", z.Metadatei.Typ, ks, es)
	}

	a.Arg = ks

	w.addCommonMatches(z, a)

	a.Text.Copy = ks
	a.Uid = "zit://" + ks

	{
		var sb strings.Builder

		if _, err := w.organizeFmt.WriteStringFormat(&sb, z); err != nil {
			a = w.errorToItem(err)
			return
		}

		a.Mods["alt"] = alfred.Mod{
			Valid:    true,
			Arg:      sb.String(),
			Subtitle: sb.String(),
		}
	}

	return
}

func (w *Writer) etikettToItem(
	z *sku.Transacted,
	e *kennung.Etikett,
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = "@" + e.String()

	a.Arg = e.String()

	w.addCommonMatches(z, a)

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
	mb.AddMatch(e.GetHead())
	mb.AddMatch(e.GetTail())

	a.Match.ReadFromBuffer(&mb.Buffer)

	a.Text.Copy = e.String()
	a.Uid = "zit://" + e.String()

	return
}
