package organize_text

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type writer struct {
	sku.ObjectFactory
	OmitLeadingEmptyLine bool
	object_metadata.Metadata
	*format.LineWriter
	maxDepth int
	options  Options
}

func (av writer) write(a *Assignment) (err error) {
	spaceCount := av.maxDepth

	hinMaxWidth := 3

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
				quiter.StringCommaSeparated(a.Transacted.Metadata.Tags),
			),
		)
		av.WriteExactlyOneEmpty()
	}

	write := func(z *obj) (err error) {
		var sb strings.Builder

		if z.tipe.IsDirectOrSelf() {
			sb.WriteString("- ")
		} else {
			sb.WriteString("% ")
		}

		cursor := z.sku.Clone()
		cursorExternal := cursor.GetSkuExternal()
		cursorExternal.Metadata.Subtract(&av.Metadata)
		mes := cursorExternal.GetMetadata().GetTags().CloneMutableSetPtrLike()

		if err = a.SubtractFromSet(mes); err != nil {
			err = errors.Wrap(err)
			return
		}

		cursorExternal.Metadata.SetTags(mes)

		if _, err = av.options.fmtBox.WriteStringFormat(
			&sb,
			cursor,
		); err != nil {
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
