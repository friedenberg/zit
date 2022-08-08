package atomic_flush

import "github.com/friedenberg/zit/alfa/errors"


//TODO WriteCloser factory for directories and files

type Bowl struct {
	Directory string
	Deposits  []Deposit
}

func (b Bowl) Flush() (err error) {
	return
}

func (b Bowl) validateMovements() (movements map[string]string, err error) {
	movements = make(map[string]string)
	oldPaths := make(map[string]Deposit)
	newPaths := make(map[string]Deposit)

	for _, d := range b.Deposits {
		if d.Error != nil {
			err = errors.Error(d.Error)
			return
		}

		if old, ok := oldPaths[d.OldPath]; ok {
			err = errors.Errorf(
				"more than one deposit with the same old path:\nfirst: %v\nsecond: %v",
				old,
				d,
			)

			return
		}

		if n, ok := newPaths[d.NewPath]; ok {
			err = errors.Errorf(
				"more than one deposit with the same new path:\nfirst: %v\nsecond: %v",
				n,
				d,
			)

			return
		}

		movements[d.OldPath] = d.NewPath
	}

	return
}
