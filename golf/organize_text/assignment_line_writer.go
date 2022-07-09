package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/bravo/line_format"
)

type assignmentLineWriter struct {
	*line_format.Writer
}

func (av assignmentLineWriter) write(a *assignment) (err error) {
	if a.etiketten.Len() > 0 {
		av.WriteLines(fmt.Sprintf("%s %s", strings.Repeat("#", a.depth), a.etiketten))
		av.WriteExactlyOneEmpty()
	}

	for z, _ := range a.named {
		av.WriteLines(fmt.Sprintf("- [%s] %s", z.Hinweis, z.Bezeichnung))
	}

	if len(a.named) > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.children {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}
