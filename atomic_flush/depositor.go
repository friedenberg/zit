package atomic_flush

import "context"

type Depositor interface {
	Flush(context.Context) Deposit
}
