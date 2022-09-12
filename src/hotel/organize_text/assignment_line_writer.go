package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/line_format"
)

type assignmentLineWriter struct {
	RightAlignedIndents bool
	*line_format.Writer
	maxDepth            int
	maxKopf, maxScwhanz int
}

func (av assignmentLineWriter) write(a *assignment) (err error) {
	if av.RightAlignedIndents {
		return av.writeRightAligned(a)
	} else {
		return av.writeNormal(a)
	}
}

func (av assignmentLineWriter) writeNormal(a *assignment) (err error) {
	tab_prefix := ""

	if a.Depth() == 0 {
		av.WriteExactlyOneEmpty()
	} else if a.Depth() < 0 {
		err = errors.Errorf("negative depth: %d", a.Depth())
		return
	} else {
		tab_prefix = strings.Repeat(" ", a.Depth()*2-(a.Depth())-1)
	}

	if a.etiketten.Len() > 0 {
		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s",
				tab_prefix,
				strings.Repeat("#", a.Depth()),
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

	if len(a.named) > 0 || len(a.unnamed) > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.children {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}

func (av assignmentLineWriter) writeRightAligned(a *assignment) (err error) {
	spaceCount := av.maxDepth

	hinMaxWidth := av.maxKopf + av.maxScwhanz + 4

	if spaceCount < hinMaxWidth {
		spaceCount = hinMaxWidth
	}

	tab_prefix := strings.Repeat(" ", hinMaxWidth)

	if a.Depth() == 0 {
		av.WriteExactlyOneEmpty()
	} else if a.Depth() < 0 {
		err = errors.Errorf("negative depth: %d", a.Depth())
		return
	}

	if a.etiketten.Len() > 0 {
		sharps := strings.Repeat("#", a.Depth())
		alignmentSpacing := strings.Repeat(" ", a.AlignmentSpacing())

		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s%s",
				tab_prefix[len(sharps)-1:],
				sharps,
				alignmentSpacing,
				a.etiketten,
			),
		)
		av.WriteExactlyOneEmpty()
	}

	for _, z := range a.unnamed.sorted() {
		av.WriteLines(
			fmt.Sprintf("- %s%s", tab_prefix, z.Bezeichnung))
	}

	for _, z := range a.named.sorted() {
		h := z.HinweisAligned(av.maxKopf, av.maxScwhanz)
		av.WriteLines(fmt.Sprintf("- [%s] %s", h, z.Bezeichnung))
	}

	if len(a.named) > 0 || len(a.unnamed) > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.children {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}
