package object_metadata_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

func AppendMetadataFields(
	fields []string_format_writer.Field,
	m *object_metadata.Metadata,
) []string_format_writer.Field {
	return append(
		fields,
		MetadataFieldSha(m),
		MetadataFieldType(m),
		MetadataFieldDescription(m),
	)
}

func MetadataFieldSha(
	m *object_metadata.Metadata,
) string_format_writer.Field {
	return string_format_writer.Field{
		Value:     "@" + m.Sha().String(),
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
