package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type PageSelectionReason int

const (
	PageSelectionReasonUnknown = PageSelectionReason(iota)
	PageSelectionReasonStoredSha
	PageSelectionReasonNamedSha
	PageSelectionReasonHinweis
	PageSelectionReasonEtikett
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

func (z *Zettel) Reset() {
	z.Transacted.Reset()
	z.PageSelection.Reason = PageSelectionReasonUnknown
	z.PageSelection.Value = ""
	z.EtikettenExpandedSorted = z.EtikettenExpandedSorted[:0]
	z.EtikettenSorted = z.EtikettenSorted[:0]
}
