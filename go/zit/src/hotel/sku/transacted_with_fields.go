package sku

import "code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"

type Field = string_format_writer.Field

type TransactedWithFields struct {
	Transacted
	Fields []Field
}
