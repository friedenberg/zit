package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/line_format"
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
		tab_prefix = strings.Repeat(" ", a.depth*2-(a.depth)-1)
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

	for _, z := range a.unnamed.sorted() {
		av.WriteLines(
			fmt.Sprintf("%s- %s", tab_prefix, z.Bezeichnung))
	}

	for _, z := range a.named.sorted() {
		av.WriteLines(fmt.Sprintf("%s- [%s] %s", tab_prefix, z.Hinweis, z.Bezeichnung))
	}

	if len(a.named) > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.childrenSorted() {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}
