package object_metadata_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

func MetadataShaString(
	m *object_metadata.Metadata,
	abbr ids.FuncAbbreviateString,
) (v string, err error) {
	s := m.Sha()
	v = s.String()

	if abbr != nil {
		var v1 string

		sh := sha.Make(s)

		if v1, err = abbr(sh); err != nil {
			err = errors.Wrap(err)
			return
		}

		if v1 != "" {
			v = v1
		} else {
			ui.Todo("abbreviate sha produced empty string")
		}
	}

	return
}

func MetadataFieldShaString(
	v string,
) string_format_writer.Field {
	return string_format_writer.Field{
		Value:     "@" + v,
		ColorType: string_format_writer.ColorTypeHash,
	}
}

func MetadataFieldType(
	m *object_metadata.Metadata,
) string_format_writer.Field {
	return string_format_writer.Field{
		Value:     "!" + m.Type.String(),
		ColorType: string_format_writer.ColorTypeType,
	}
}

func MetadataFieldDescription(
	m *object_metadata.Metadata,
) string_format_writer.Field {
	return string_format_writer.Field{
		Value:     m.Description.String(),
		ColorType: string_format_writer.ColorTypeUserData,
	}
}