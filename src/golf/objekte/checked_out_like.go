package objekte

type CheckedOutLike interface {
	GetInternal() TransactedLike // TODO-P0 rename to GetInternalLike
	GetExternal() ExternalLike   // TODO-P0 rename to GetExternalLike
	GetState() CheckedOutState
}
