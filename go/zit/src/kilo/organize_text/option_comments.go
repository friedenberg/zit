package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type optionComment interface {
	GetOptionComment() optionComment
	schnittstellen.Stringer
	ApplyToText(Options, *Text) error
	ApplyToReader(Options, *assignmentLineReader) error
	ApplyToWriter(Options, *assignmentLineWriter) error
}

type optionCommentFactory struct{}

func (ocf optionCommentFactory) Make(c string) (oc optionComment, err error) {
	c = strings.TrimSpace(c)

	head, tail, found := strings.Cut(c, ":")

	if !found {
		// err = errors.New("':' not found")
		return
	}

	switch head {
	case "format":
		oc = optionCommentFormat(tail)

	case "hide":
		oc = optionCommentHide(tail)

	default:
	}

	return
}

type optionCommentFormat string

func (ocf optionCommentFormat) GetOptionComment() optionComment {
	return ocf
}

func (ocf optionCommentFormat) String() string {
	return fmt.Sprintf("format:%s", string(ocf))
}

func (ocf optionCommentFormat) ApplyToText(Options, *Text) (err error) {
	return
}

func (ocf optionCommentFormat) ApplyToReader(
	Options,
	*assignmentLineReader,
) (err error) {
	return
}

func (ocf optionCommentFormat) ApplyToWriter(
	f Options,
	aw *assignmentLineWriter,
) (err error) {
	switch string(ocf) {
	case "new":
		aw.stringFormatWriter = &f.organizeNew

	case "old":
		aw.stringFormatWriter = &f.organize

	default:
		err = collections.MakeErrNotFoundString(string(ocf))
		return
	}

	return
}

type optionCommentHide string

func (ocf optionCommentHide) GetOptionComment() optionComment {
	return ocf
}

func (ocf optionCommentHide) String() string {
	return fmt.Sprintf("hide:%s", string(ocf))
}

func (ocf optionCommentHide) ApplyToText(Options, *Text) (err error) {
	return
}

func (ocf optionCommentHide) ApplyToReader(
	Options,
	*assignmentLineReader,
) (err error) {
	return
}

func (ocf optionCommentHide) ApplyToWriter(
	f Options,
	aw *assignmentLineWriter,
) (err error) {
	return
}
