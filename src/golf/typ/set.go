package typ

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Set = collections.ValueSet[kennung.Typ, *kennung.Typ]
