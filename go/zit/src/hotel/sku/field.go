package sku

import "code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"

type FieldType int

const (
	FieldTypeUnknown = FieldType(iota)
	FieldTypeId
	FieldTypeHash
	FieldTypeTag
	FieldTypeType
	FieldTypeUserData
)

type Field struct {
	FieldType
	Value interfaces.Stringer
}
