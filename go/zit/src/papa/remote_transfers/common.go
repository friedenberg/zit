package remote_transfers

import (
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type common struct {
	*umwelt.Umwelt
	pmf object_inventory_format.Format
}
