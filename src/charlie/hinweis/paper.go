package hinweis

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Abbr interface {
	AbbreviateHinweis(h Hinweis) (ha Hinweis, err error)
}

type Paper struct {
	MaxKopf, MaxSchwanz int
	paper
}

type paper struct {
	Abbr
	Hinweis
}

func MakePaper(h Hinweis, a Abbr) (p *Paper) {
	return &Paper{
		paper: paper{
			Abbr:    a,
			Hinweis: h,
		},
	}
}

func (p Paper) WriteTo(w io.Writer) (n int64, err error) {
	h := p.Hinweis

	if p.Abbr != nil {
		if h, err = p.AbbreviateHinweis(p.Hinweis); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	parts := strings.Split(h.String(), "/")

	diffKopf := p.MaxKopf - len(parts[0])
	if diffKopf > 0 {
		parts[0] = strings.Repeat(" ", diffKopf) + parts[0]
	}

	diffSchwanz := p.MaxSchwanz - len(parts[1])
	if diffSchwanz > 0 {
		parts[1] = parts[1] + strings.Repeat(" ", diffSchwanz)
	}

	if _, err = io.WriteString(w, fmt.Sprintf("%s/%s", parts[0], parts[1])); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
