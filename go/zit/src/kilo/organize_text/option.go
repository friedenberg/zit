package organize_text

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
)

type Option interface {
	GetOption() Option
	interfaces.Stringer
	ApplyToText(Options, *Assignment) error
	ApplyToReader(Options, *assignmentLineReader) error
	ApplyToWriter(Options, *assignmentLineWriter) error
}

type optionCommentFactory struct{}

func (ocf optionCommentFactory) Make(c string) (oc Option, err error) {
	c = strings.TrimSpace(c)

	head, tail, found := strings.Cut(c, ":")

	if !found {
		if c == "abort" {
			err = errors.New("aborting!")
			return
		}
		// err = errors.New("':' not found")
		return
	}

	switch head {
	case "format":
		oc = optionCommentFormat(tail)

	case "hide":
		oc = optionCommentHide(tail)

	case "dry-run":
		var b values.Bool

		if err = b.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

		oc = optionCommentDryRun(b)

	default:
	}

	return
}

type optionCommentFormat string

func (ocf optionCommentFormat) GetOption() Option {
	return ocf
}

func (ocf optionCommentFormat) String() string {
	return fmt.Sprintf("format:%s", string(ocf))
}

func (ocf optionCommentFormat) ApplyToText(Options, *Assignment) (err error) {
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
		aw.stringFormatWriter = f.stringFormatWriter

	default:
		err = collections.MakeErrNotFoundString(string(ocf))
		return
	}

	return
}

type optionCommentHide string

func (ocf optionCommentHide) GetOption() Option {
	return ocf
}

func (ocf optionCommentHide) String() string {
	return fmt.Sprintf("hide:%s", string(ocf))
}

func (ocf optionCommentHide) ApplyToText(Options, *Assignment) (err error) {
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

type optionCommentDryRun values.Bool

func (ocf optionCommentDryRun) GetOption() Option {
	return ocf
}

func (ocf optionCommentDryRun) String() string {
	return fmt.Sprintf("dry-run:%s", values.Bool(ocf))
}

func (ocf optionCommentDryRun) ApplyToText(o Options, a *Assignment) (err error) {
	o.Config.DryRun = values.Bool(ocf).Bool()
	return
}

func (ocf optionCommentDryRun) ApplyToReader(
	o Options,
	a *assignmentLineReader,
) (err error) {
	o.Config.DryRun = values.Bool(ocf).Bool()
	return
}

func (ocf optionCommentDryRun) ApplyToWriter(
	f Options,
	aw *assignmentLineWriter,
) (err error) {
	f.Config.DryRun = values.Bool(ocf).Bool()
	return
}
