package remote_transfers

import (
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type common struct {
	*umwelt.Umwelt
	pmf objekte_format.Format
}
