package sha

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Abbr interface {
	AbbreviateSha(Sha) (string, error)
}

type Paper struct {
	paper
}

type paper struct {
	Abbr
	Sha
}

func MakePaper(s Sha, a Abbr) (p *Paper) {
	return &Paper{
		paper: paper{
			Abbr: a,
			Sha:  s,
		},
	}
}

func (p paper) WriteTo(w io.Writer) (n int64, err error) {
	var sha string

	if sha, err = p.AbbreviateSha(p.Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.WriteString(w, sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
