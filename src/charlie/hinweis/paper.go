package hinweis

import (
	"io"

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

	if _, err = io.WriteString(w, h.Aligned(p.MaxKopf, p.MaxSchwanz)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
