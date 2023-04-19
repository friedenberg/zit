package objekte

type CheckedOutLike interface {
	GetInternalLike() TransactedLike
	GetExternalLike() ExternalLike
	GetState() CheckedOutState
}

type CheckedOutLikePtr interface {
	CheckedOutLike
	DetermineState()
}
