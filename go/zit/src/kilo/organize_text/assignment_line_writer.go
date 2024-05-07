package organize_text

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/format"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
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

func (av assignmentLineWriter) write(a *Assignment) (err error) {
	if av.RightAlignedIndents {
		return av.writeRightAligned(a)
	} else {
		return av.writeNormal(a)
	}
}

func (av assignmentLineWriter) writeNormal(a *Assignment) (err error) {
	tab_prefix := ""

	if a.GetDepth() == 0 && !av.OmitLeadingEmptyLine {
		av.WriteExactlyOneEmpty()
	} else if a.GetDepth() < 0 {
		err = errors.Errorf("negative depth: %d", a.GetDepth())
		return
	} else {
		tab_prefix = strings.Repeat(" ", a.GetDepth()*2-(a.GetDepth())-1)
	}

	if a.Etiketten.Len() > 0 {
		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s",
				tab_prefix,
				strings.Repeat("#", a.GetDepth()),
				iter.StringCommaSeparated(a.Etiketten),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	for _, z := range sortObjSet(a.Unnamed) {
		av.WriteLines(
			fmt.Sprintf("%s- %s", tab_prefix, z.Metadatei.Bezeichnung),
		)
	}

	cursor := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(cursor)

	for _, z := range sortObjSet(a.Named) {
		var sb strings.Builder

		sb.WriteString(tab_prefix)
		sb.WriteString("- ")

		sku.TransactedResetter.ResetWith(cursor, &z.Transacted)
		cursor.Metadatei.Subtract(&av.Metadatei)

		if _, err = av.stringFormatWriter.WriteStringFormat(&sb, cursor); err != nil {
			err = errors.Wrap(err)
			return
		}

		av.WriteStringers(&sb)
	}

	if a.Named.Len() > 0 || a.Unnamed.Len() > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.Children {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}

func (av assignmentLineWriter) writeRightAligned(a *Assignment) (err error) {
	spaceCount := av.maxDepth

	kopfUndSchwanz := av.maxKopf + av.maxSchwanz

	hinMaxWidth := 4

	hinMaxWidth += kopfUndSchwanz - 1

	if spaceCount < hinMaxWidth {
		spaceCount = hinMaxWidth
	}

	tab_prefix := strings.Repeat(" ", spaceCount+1)

	if a.GetDepth() == 0 && !av.OmitLeadingEmptyLine {
		av.WriteExactlyOneEmpty()
	} else if a.GetDepth() < 0 {
		err = errors.Errorf("negative depth: %d", a.GetDepth())
		return
	}

	if a.Etiketten != nil && a.Etiketten.Len() > 0 {
		sharps := strings.Repeat("#", a.GetDepth())
		alignmentSpacing := strings.Repeat(" ", a.AlignmentSpacing())

		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s%s",
				tab_prefix[len(sharps)-1:],
				sharps,
				alignmentSpacing,
				iter.StringCommaSeparated(a.Etiketten),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	cursor := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(cursor)

	write := func(z *sku.Transacted) (err error) {
		var sb strings.Builder

		sb.WriteString("- ")
		sku.TransactedResetter.ResetWith(cursor, z)
		cursor.Metadatei.Subtract(&av.Metadatei)

		if err = a.SubtractFromSet(
			cursor.Metadatei.GetEtikettenMutable(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = av.stringFormatWriter.WriteStringFormat(&sb, cursor); err != nil {
			err = errors.Wrap(err)
			return
		}

		av.WriteStringers(&sb)

		return
	}

	for _, z := range sortObjSet(a.Unnamed) {
		if err = write(&z.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, z := range sortObjSet(a.Named) {
		if err = write(&z.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if a.Named.Len() > 0 || a.Unnamed.Len() > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.Children {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}
