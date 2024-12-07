package remote_conn

// TODO-P4 rename to RemoteRequest
//
//go:generate stringer -type=DialogueType
type DialogueType int

const (
	DialogueTypeMain = DialogueType(iota)
	DialogueTypeSkusForFilter
	DialogueTypeObjects
	DialogueTypeBlobs
	DialogueTypeObjectReader
	DialogueTypeBlobReader
	DialogueTypeObjectWriter
	DialogueTypeBlobWriter
	DialogueTypePull
	DialogueTypePullBLob
	DialogueTypePush
	DialogueTypePushObjects
	DialogueTypePushBlob
	DialogueTypeGetNeededSkus
)
