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
	*errors.Ctx
}

type paper struct {
	*strings.Builder
}

func Make(f *os.File, ctx *errors.Ctx) (pa *Paper) {
	pa = &Paper{
		paper: paper{
			Builder: &strings.Builder{},
		},
		file: f,
		Ctx:  ctx,
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

func (p *Paper) Print() {
	p.NewLine()
	_, p.Err = io.WriteString(p.file, p.String())
}
