package changes

import (
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

type ChangeBezeichnungKeyer struct{}

func (ChangeBezeichnungKeyer) GetKey(c ChangeBezeichnung) string {
	return c.Kennung
}

func (ChangeBezeichnungKeyer) GetKeyPtr(c *ChangeBezeichnung) string {
	return c.Kennung
}

type ChangeBezeichnung struct {
	Kennung     string
	Bezeichnung bezeichnung.Bezeichnung
}
