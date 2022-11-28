package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
	"github.com/friedenberg/zit/src/typ_toml"
)

type Objekte = typ_toml.Objekte
type Transacted = objekte.Transacted2[typ_toml.Objekte, *typ_toml.Objekte, kennung.Typ, *kennung.Typ]
