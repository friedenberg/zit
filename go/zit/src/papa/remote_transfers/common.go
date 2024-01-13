package remote_transfers

import (
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type common struct {
	*umwelt.Umwelt
	pmf objekte_format.Format
}
