package typ

import (
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
	"github.com/friedenberg/zit/src/foxtrot/typ_toml"
)

type Objekte = typ_toml.Objekte
type Transacted = objekte.Transacted[typ_toml.Objekte, *typ_toml.Objekte, kennung.Typ, *kennung.Typ]
