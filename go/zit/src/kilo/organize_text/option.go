package organize_text

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
)

type OptionComment interface {
	CloneOptionComment() OptionComment
	interfaces.StringerSetter
}

type OptionCommentWithApply interface {
	OptionComment
	ApplyToText(Options, *Assignment) error
	ApplyToReader(Options, *reader) error
	ApplyToWriter(Options, *writer) error
}

// TODO add config to automatically add dry run if necessary
func MakeOptionCommentSet(
	elements map[string]OptionComment,
	options ...OptionComment,
) OptionCommentSet {
	ocs := OptionCommentSet{
		prototype:      make(PrototypeOptionComments),
		OptionComments: options,
	}

	if elements != nil {
		for k, el := range elements {
			ocs.AddPrototype(k, el)
		}
	}

	ocs.AddPrototype("format", optionCommentFormat(""))
	ocs.AddPrototype("hide", optionCommentHide(""))
	ocs.AddPrototype("", optionCommentHide(""))

	return ocs
}

type PrototypeOptionComments map[string]OptionComment

type OptionCommentSet struct {
	prototype      PrototypeOptionComments
	OptionComments []OptionComment
}

func (ocs *OptionCommentSet) GetPrototypeOptionComments() PrototypeOptionComments {
	return ocs.prototype
}

func (ocs *OptionCommentSet) AddPrototype(
	key string,
	o OptionComment,
) OptionComment {
	o = OptionCommentWithKey{
		Key:           key,
		OptionComment: o,
	}

	ocs.prototype[key] = o

	return o
}

func (ocs *OptionCommentSet) AddPrototypeAndOption(
	key string,
	o OptionComment,
) OptionComment {
	o = ocs.AddPrototype(key, o)
	ocs.OptionComments = append(ocs.OptionComments, o)
	return o
}

func (ocs *OptionCommentSet) Set(v string) (err error) {
	head, tail, _ := strings.Cut(v, ":")

	oc, ok := ocs.prototype[head]

	if ok {
		oc = oc.CloneOptionComment()
	} else {
		oc = OptionCommentUnknown("")
	}

	oc = OptionCommentWithKey{
		Key:           head,
		OptionComment: oc,
	}

	if err = oc.Set(tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	ocs.OptionComments = append(
		ocs.OptionComments,
		oc,
	)

	return
}

// TODO add support for ApplyTo*
type OptionCommentWithKey struct {
	Key string
	OptionComment
}

func (ocf OptionCommentWithKey) CloneOptionComment() OptionComment {
	return ocf
}

func (ocf OptionCommentWithKey) Set(v string) (err error) {
	if err = ocf.OptionComment.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ocf OptionCommentWithKey) String() string {
	return fmt.Sprintf("%s:%s", ocf.Key, ocf.OptionComment)
}

type optionCommentFormat string

func (ocf optionCommentFormat) CloneOptionComment() OptionComment {
	return ocf
}

func (ocf optionCommentFormat) Set(v string) (err error) {
	return todo.Implement()
}

func (ocf optionCommentFormat) String() string {
	return fmt.Sprintf("format:%s", string(ocf))
}

func (ocf optionCommentFormat) ApplyToText(Options, *Assignment) (err error) {
	return
}

func (ocf optionCommentFormat) ApplyToReader(
	Options,
	*reader,
) (err error) {
	return
}

func (ocf optionCommentFormat) ApplyToWriter(
	f Options,
	aw *writer,
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

func (ocf optionCommentHide) CloneOptionComment() OptionComment {
	return ocf
}

func (ocf optionCommentHide) Set(v string) (err error) {
	return todo.Implement()
}

func (ocf optionCommentHide) String() string {
	return fmt.Sprintf("hide:%s", string(ocf))
}

func (ocf optionCommentHide) ApplyToText(Options, *Assignment) (err error) {
	return
}

func (ocf optionCommentHide) ApplyToReader(
	Options,
	*reader,
) (err error) {
	return
}

func (ocf optionCommentHide) ApplyToWriter(
	f Options,
	aw *writer,
) (err error) {
	return
}

type OptionCommentDryRun struct {
	interfaces.ConfigDryRun
}

func (ocf *OptionCommentDryRun) CloneOptionComment() OptionComment {
	return ocf
}

func (ocf *OptionCommentDryRun) Set(v string) (err error) {
	var b values.Bool

	if err = b.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	ocf.SetDryRun(b.Bool())

	return
}

func (ocf *OptionCommentDryRun) String() string {
	return fmt.Sprintf("%t", ocf.IsDryRun())
}

type OptionCommentUnknown string

func (ocf OptionCommentUnknown) CloneOptionComment() OptionComment {
	return ocf
}

func (ocf OptionCommentUnknown) Set(v string) (err error) {
	ocf = OptionCommentUnknown(strings.TrimSpace(v))
	return
}

func (ocf OptionCommentUnknown) String() string {
	return string(ocf)
}

type OptionCommentBooleanFlag struct {
	Value   *bool
	Comment string
}

func (ocf OptionCommentBooleanFlag) CloneOptionComment() OptionComment {
	return ocf
}

func (ocf OptionCommentBooleanFlag) Set(v string) (err error) {
	head, tail, _ := strings.Cut(v, " ")

	var b values.Bool

	if err = b.Set(head); err != nil {
		err = errors.Wrap(err)
		return
	}

	*ocf.Value = b.Bool()

	ocf.Comment = tail

	return
}

func (ocf OptionCommentBooleanFlag) String() string {
	if ocf.Comment != "" {
		return fmt.Sprintf("%t %s", *ocf.Value, ocf.Comment)
	} else {
		return fmt.Sprintf("%t", *ocf.Value)
	}
}
