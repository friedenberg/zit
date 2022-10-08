package store_verzeichnisse

import "io"

type PageSelectionReason int

const (
	PageSelectionReasonUnknown = PageSelectionReason(iota)
	PageSelectionReasonStoredSha
	PageSelectionReasonNamedSha //unused
	PageSelectionReasonHinweis
	PageSelectionReasonEtikett
	PageSelectionReasonAkte
)

func (psr PageSelectionReason) WriteZettelVerzeichnisse(
	z *Zettel,
) (err error) {
	if z.PageSelection.Reason != psr {
		err = io.EOF
	}

	return
}
