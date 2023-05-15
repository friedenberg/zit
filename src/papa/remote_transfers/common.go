package remote_transfers

import (
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type common struct {
	*umwelt.Umwelt
	pmf persisted_metadatei_format.Format
}
