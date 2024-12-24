package errors

type Helpful interface {
	error
	GetHelpfulError() Helpful
	ErrorCause() []string
	ErrorRecovery() []string
}
