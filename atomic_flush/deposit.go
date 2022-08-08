package atomic_flush

type Deposit struct {
	OldPath, NewPath string
	Error            error
}
