package metadatei

import (
	"fmt"

	"github.com/friedenberg/zit/src/delta/kennung"
)

type ErrHasInlineAkteAndFilePath struct {
	Metadatei Metadatei
	AkteFD    kennung.FD
}

func (e ErrHasInlineAkteAndFilePath) Error() string {
	return fmt.Sprintf(
		"zettel text has both inline akte and filepath: \nexternal path: %s\nexternal sha: %s\ninline sha: %s",
		e.AkteFD.Path,
		e.AkteFD.Sha,
		e.Metadatei.AkteSha,
	)
}
