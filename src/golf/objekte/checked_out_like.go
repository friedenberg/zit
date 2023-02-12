package objekte

type CheckedOutLike interface {
	GetInternal() TransactedLike
	GetExternal() ExternalLike
}
