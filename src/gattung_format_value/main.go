package gattung_format_value

type Format int

const (
	FormatUnknown = Format(iota)
	FormatAkte
	FormatJson
	FormatLog
	FormatObjekte
	FormatText
	FormatToml
)
