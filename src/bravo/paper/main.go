package paper

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Paper struct {
	paper
	file *os.File
}

type paper struct {
	*strings.Builder
}

func Make(f *os.File) (pa *Paper, err error) {
	pa = &Paper{
		paper: paper{
			Builder: &strings.Builder{},
		},
		file: f,
	}

	return
}

func (p *Paper) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, strings.NewReader(p.String()))
}

func (p *Paper) WriteFrom(ws ...io.WriterTo) (err error) {
	for _, w := range ws {
		if _, err = w.WriteTo(p.Builder); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (p *Paper) WriteFormat(f string, vs ...interface{}) {
	s := fmt.Sprintf(f, vs...)
	p.WriteString(s)
}

func (p *Paper) NewLine() {
	p.WriteString("\n")
}

func (p *Paper) Print() (err error) {
	p.NewLine()
	_, err = io.WriteString(p.file, p.String())

	return
}
