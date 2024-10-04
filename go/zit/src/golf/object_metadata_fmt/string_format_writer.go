package object_metadata_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type Box struct {
	options       print_options.General
	contentPrefix string

	writeTyp         bool
	writeDescription bool
	writeEtiketten   bool

	Sha    interfaces.StringFormatWriter[interfaces.Sha]
	Type   interfaces.StringFormatWriter[*ids.Type]
	Fields interfaces.StringFormatWriter[string_format_writer.Box]
	Tags   interfaces.StringFormatWriter[*ids.Tag]
}

func MakeBoxMetadataFormat(
	options print_options.General,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	typStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
	etikettenStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *Box {
	return &Box{
		options: options,
		contentPrefix: string_format_writer.StringPrefixFromOptions(
			options,
		),
		writeTyp:         true,
		writeDescription: true,
		writeEtiketten:   true,
		Sha:              shaStringFormatWriter,
		Type:             typStringFormatWriter,
		Fields:           fieldsFormatWriter,
		Tags:             etikettenStringFormatWriter,
	}
}
