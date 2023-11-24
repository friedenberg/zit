package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type assignmentLineWriter struct {
	RightAlignedIndents  bool
	OmitLeadingEmptyLine bool
	Metadatei            metadatei.Metadatei
	*format.LineWriter
	maxDepth            int
	maxKopf, maxSchwanz int
	maxLen              int
	stringFormatWriter  schnittstellen.StringFormatWriter[*sku.Transacted]
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
				iter.StringCommaSeparated[kennung.Etikett](a.etiketten),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	for _, z := range sortObjSet(a.unnamed) {
		av.WriteLines(
			fmt.Sprintf("%s- %s", tab_prefix, z.Sku.Metadatei.Bezeichnung),
		)
	}

	cursor := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(cursor)

	for _, z := range sortObjSet(a.named) {
		var sb strings.Builder

		sb.WriteString(tab_prefix)
		sb.WriteString("- ")

		sku.TransactedResetter.ResetWithPtr(cursor, &z.Sku)
		cursor.Metadatei.Subtract(&av.Metadatei)

		if _, err = av.stringFormatWriter.WriteStringFormat(&sb, cursor); err != nil {
			err = errors.Wrap(err)
			return
		}

		av.WriteStringers(&sb)
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

	kopfUndSchwanz := av.maxKopf + av.maxSchwanz

	hinMaxWidth := 4

	hinMaxWidth += kopfUndSchwanz - 1

	if spaceCount < hinMaxWidth {
		spaceCount = hinMaxWidth
	}

	tab_prefix := strings.Repeat(" ", spaceCount+1)

	if a.Depth() == 0 && !av.OmitLeadingEmptyLine {
		av.WriteExactlyOneEmpty()
	} else if a.Depth() < 0 {
		err = errors.Errorf("negative depth: %d", a.Depth())
		return
	}

	if a.etiketten != nil && a.etiketten.Len() > 0 {
		sharps := strings.Repeat("#", a.Depth())
		alignmentSpacing := strings.Repeat(" ", a.AlignmentSpacing())

		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s%s",
				tab_prefix[len(sharps)-1:],
				sharps,
				alignmentSpacing,
				iter.StringCommaSeparated[kennung.Etikett](a.etiketten),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	for _, z := range sortObjSet(a.unnamed) {
		av.WriteLines(
			fmt.Sprintf("- %s%s", tab_prefix, z.Sku.Metadatei.Bezeichnung),
		)
	}

	cursor := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(cursor)

	for _, z := range sortObjSet(a.named) {
		var sb strings.Builder

		sb.WriteString("- ")
		sku.TransactedResetter.ResetWithPtr(cursor, &z.Sku)
		cursor.Metadatei.Subtract(&av.Metadatei)

		if _, err = av.stringFormatWriter.WriteStringFormat(&sb, cursor); err != nil {
			err = errors.Wrap(err)
			return
		}

		av.WriteStringers(&sb)
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
