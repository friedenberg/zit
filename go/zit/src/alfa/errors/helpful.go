package errors

type Helpful interface {
	error
	GetHelpfulError() Helpful
	ErrorCause() []string
	ErrorRecovery() []string
}

type Retryable interface {
	GetRetryableError() Retryable
	Recover(IContext)
}
