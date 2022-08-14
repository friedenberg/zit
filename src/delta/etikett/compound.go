package etikett

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
)

type Compound []Etikett

func (c Compound) Strings() (out []string) {
	out = make([]string, len(c))

	for i, _ := range out {
		out[i] = c[i].String()
	}

	return
}

func (c Compound) String() string {
	return strings.Join(c.Strings(), ",")
}

func (c *Compound) Set(v string) (err error) {
	es := strings.Split(v, ",")
	*c = make([]Etikett, len(es))

	for i, e := range *c {
		if err = e.Set(es[i]); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
