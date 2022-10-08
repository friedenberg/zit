package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type PageSelection struct {
	Reason PageSelectionReason
	Value  string
}

type Zettel struct {
	PageSelection           PageSelection
	Transacted              zettel_transacted.Zettel
	EtikettenExpandedSorted []string
	EtikettenSorted         []string
}

func (i *Zettelen) MakeZettel(
	tz zettel_transacted.Zettel,
	reason PageSelectionReason,
	value string,
) (z *Zettel) {
	z = i.pool.Get()
	z.PageSelection.Reason = reason
	z.PageSelection.Value = value
	z.Transacted = tz
	z.EtikettenExpandedSorted = tz.Named.Stored.Zettel.Etiketten.Expanded().SortedString()
	z.EtikettenSorted = tz.Named.Stored.Zettel.Etiketten.SortedString()

	return
}

func (z *Zettel) Reset() {
	z.Transacted.Reset()
	z.PageSelection.Reason = PageSelectionReasonUnknown
	z.PageSelection.Value = ""
	z.EtikettenExpandedSorted = z.EtikettenExpandedSorted[:0]
	z.EtikettenSorted = z.EtikettenSorted[:0]
}
