package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
)

type assignmentLineWriter struct {
	RightAlignedIndents  bool
	OmitLeadingEmptyLine bool
	*format.LineWriter
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

	if a.Depth() == 0 && !av.OmitLeadingEmptyLine {
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
				collections.StringCommaSeparated(a.etiketten),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	for _, z := range sortNewZettelSet(a.unnamed) {
		av.WriteLines(
			fmt.Sprintf("%s- %s", tab_prefix, z.Bezeichnung))
	}

	for _, z := range sortZettelSet(a.named) {
		av.WriteLines(fmt.Sprintf("%s- [%s] %s", tab_prefix, z.Hinweis, z.Bezeichnung))
	}

	if a.named.Len() > 0 || a.unnamed.Len() > 0 {
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

	if a.Depth() == 0 && !av.OmitLeadingEmptyLine {
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
				collections.StringCommaSeparated(a.etiketten),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	for _, z := range sortNewZettelSet(a.unnamed) {
		av.WriteLines(
			fmt.Sprintf("- %s%s", tab_prefix, z.Bezeichnung))
	}

	for _, z := range sortZettelSet(a.named) {
		h := z.Hinweis.Aligned(av.maxKopf, av.maxScwhanz)
		av.WriteLines(fmt.Sprintf("- [%s] %s", h, z.Bezeichnung))
	}

	if a.named.Len() > 0 || a.unnamed.Len() > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.children {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}
