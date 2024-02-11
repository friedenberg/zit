package changes

import (
	"code.linenisgreat.com/zit-go/src/echo/bezeichnung"
)

type ChangeBezeichnungKeyer struct{}

func (ChangeBezeichnungKeyer) GetKey(c *ChangeBezeichnung) string {
	return c.Kennung
}

type ChangeBezeichnung struct {
	Kennung     string
	Bezeichnung bezeichnung.Bezeichnung
}
