package organize_text

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type assignmentLineWriter struct {
	RightAlignedIndents  bool
	OmitLeadingEmptyLine bool
	object_metadata.Metadata
	*format.LineWriter
	maxDepth           int
	maxHead, maxTail   int
	maxLen             int
	stringFormatWriter interfaces.StringFormatWriter[sku.ExternalLike]
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

	if a.Transacted.Metadata.Tags.Len() > 0 {
		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s",
				tab_prefix,
				strings.Repeat("#", a.GetDepth()),
				iter.StringCommaSeparated(a.Transacted.Metadata.Tags),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	cursor := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(cursor)

	a.Objects.Sort()

	for _, z := range a.Objects {
		var sb strings.Builder

		sb.WriteString(tab_prefix)

		if z.IsDirectOrSelf() {
			sb.WriteString("- ")
		} else {
			sb.WriteString("% ")
		}

		sku.TransactedResetter.ResetWith(cursor, z.ExternalLike.GetSku())
		cursor.Metadata.Subtract(&av.Metadata)

		if _, err = av.stringFormatWriter.WriteStringFormat(&sb, cursor); err != nil {
			err = errors.Wrap(err)
			return
		}

		av.WriteStringers(&sb)
	}

	if a.Objects.Len() > 0 {
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

	kopfUndSchwanz := av.maxHead + av.maxTail

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

	if a.Transacted.Metadata.Tags != nil && a.Transacted.Metadata.Tags.Len() > 0 {
		sharps := strings.Repeat("#", a.GetDepth())
		alignmentSpacing := strings.Repeat(" ", a.AlignmentSpacing())

		av.WriteLines(
			fmt.Sprintf(
				"%s%s %s%s",
				tab_prefix[len(sharps)-1:],
				sharps,
				alignmentSpacing,
				iter.StringCommaSeparated(a.Transacted.Metadata.Tags),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	write := func(z *obj) (err error) {
		var sb strings.Builder

		if z.IsDirectOrSelf() {
			sb.WriteString("- ")
		} else {
			sb.WriteString("% ")
		}

		cursor := z.ExternalLike.Clone()
		sk := cursor.GetSku()
		sk.Metadata.Subtract(&av.Metadata)
		mes := sk.GetMetadata().GetTags().CloneMutableSetPtrLike()

		if err = a.SubtractFromSet(mes); err != nil {
			err = errors.Wrap(err)
			return
		}

		sk.Metadata.SetTags(mes)

		if _, err = av.stringFormatWriter.WriteStringFormat(&sb, cursor); err != nil {
			err = errors.Wrap(err)
			return
		}

		av.WriteStringers(&sb)

		return
	}

	a.Objects.Sort()

	for _, z := range a.Objects {
		if err = write(z); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if a.Objects.Len() > 0 {
		av.WriteExactlyOneEmpty()
	}

	for _, c := range a.Children {
		av.write(c)
	}

	av.WriteExactlyOneEmpty()

	return
}
