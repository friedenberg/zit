package log

import (
	"log"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

var SetOutput = log.SetOutput

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

func DebugAllowCommit() errors.DevPrinter {
	return errors.Debug()
}
