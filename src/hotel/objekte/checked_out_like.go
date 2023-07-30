package objekte

type CheckedOutLike interface {
	GetInternalLike() TransactedLikePtr
	GetExternalLike() ExternalLike
	GetState() CheckedOutState
}

type CheckedOutLikePtr interface {
	CheckedOutLike
	GetExternalLikePtr() ExternalLikePtr
	DetermineState(justCheckedOut bool)
}
