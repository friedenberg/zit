package changes2

import (
	"code.linenisgreat.com/zit/src/echo/bezeichnung"
)

type ChangeBezeichnungKeyer struct{}

func (ChangeBezeichnungKeyer) GetKey(c *ChangeBezeichnung) string {
	return c.Kennung
}

type ChangeBezeichnung struct {
	Kennung     string
	Bezeichnung bezeichnung.Bezeichnung
}
