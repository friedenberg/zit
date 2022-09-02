package hinweis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type HinweisWithIndex struct {
	Hinweis
	Index int
}

func (h HinweisWithIndex) String() string {
	if h.Index < 0 {
		return fmt.Sprintf("%s", h.Hinweis)
	} else {
		return fmt.Sprintf("%s@%d", h.Hinweis, h.Index)
	}
}

func (h *HinweisWithIndex) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	vs := strings.Split(v, "@")

	switch len(vs) {

	case 1:
		h.Index = -1
		return h.Hinweis.Set(v)

	default:
		err = errors.Errorf("wrong format for HinweisWithIndex: %s", v)
		return

	case 2:
		break
	}

	if err = h.Hinweis.Set(vs[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	if h.Index, err = strconv.Atoi(vs[1]); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
