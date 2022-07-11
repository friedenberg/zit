package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/line_format"
)

type assignmentLineWriter struct {
	*line_format.Writer
}

func (av assignmentLineWriter) write(a *assignment) (err error) {
	tab_prefix := ""

	if a.depth == 0 {
		av.WriteExactlyOneEmpty()
	} else if a.depth < 0 {
		err = errors.Errorf("negative depth: %d", a.depth)
		return
	} else {
		//a.depty > 1
		tab_prefix = strings.Repeat("\t", a.depth-1)
	}

	if a.etiketten.Len() > 0 {
		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s",
				tab_prefix,
				strings.Repeat("#", a.depth),
				a.etiketten,
			),
		)
		av.WriteExactlyOneEmpty()
	}

	for z, _ := range a.named {
		av.WriteLines(
			fmt.Sprintf("%s- [%s] %s", tab_prefix,
				z.Hinweis, z.Bezeichnung))
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
