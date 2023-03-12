package log

import "github.com/friedenberg/zit/src/alfa/errors"

func Out() errors.ProdPrinter {
	return errors.Out()
}

func Err() errors.ProdPrinter {
	return errors.Err()
}

func Log() errors.DevPrinter {
	return errors.Log()
}

func Debug() errors.DevPrinter {
	return errors.Debug()
}
