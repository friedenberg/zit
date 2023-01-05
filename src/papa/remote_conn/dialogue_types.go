package remote_conn

// TODO-P4 rename to RemoteRequest
//
//go:generate stringer -type=DialogueType
type DialogueType int

const (
	DialogueTypeMain = DialogueType(iota)
	DialogueTypeSkusForFilter
	DialogueTypeObjekten
	DialogueTypeAkten
	DialogueTypeObjekteReaderForSku
	DialogueTypeAkteReaderForSha
	DialogueTypePull
	DialogueTypePullAkte
	DialogueTypePush
	DialogueTypePushObjekten
	DialogueTypePushAkte
)
