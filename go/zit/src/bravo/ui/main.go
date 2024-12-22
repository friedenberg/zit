package ui

import (
	"log"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

var (
	SetOutput = log.SetOutput
	verbose   bool
	isTest    bool
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

func SetVerbose(on bool) {
	printerLog.on = on
	printerDebug.on = on
	verbose = on

	if on {
		log.Print("verbose")
	}
}

func SetTesting() {
	isTest = true
	errors.SetTesting()
	SetVerbose(true)
}

func IsVerbose() bool {
	return verbose
}

type ProdPrinter interface {
	Print(v ...interface{}) error
	Printf(format string, v ...interface{}) error
}

type DevPrinter interface {
	ProdPrinter
	Caller(i int, vs ...interface{})
	FunctionName(skip int)
}

var (
	printerOut, printerErr   prodPrinter
	printerLog, printerDebug devPrinter
	printerBatsTestBody      devPrinter
)

func init() {
	printerOut = prodPrinter{
		f:  os.Stdout,
		on: true,
	}

	printerErr = prodPrinter{
		f:  os.Stderr,
		on: true,
	}

	printerLog = devPrinter{
		prodPrinter: prodPrinter{
			f: os.Stderr,
		},
		includesStack: true,
	}

	printerDebug = devPrinter{
		prodPrinter: prodPrinter{
			f: os.Stderr,
			// TODO-P2 determine thru compilation
			on: true,
		},
		includesStack: true,
	}

	// TODO-P2 determine thru compilation
	envVarFilter := "BATS_TEST_BODY"
	_, printerBatsTestBodyOn := os.LookupEnv(envVarFilter)

	printerBatsTestBody = devPrinter{
		prodPrinter: prodPrinter{
			f: os.Stderr,
			// TODO-P2 determine thru compilation
			on: printerBatsTestBodyOn,
		},
		includesStack: true,
	}
}

func Out() ProdPrinter {
	return printerOut
}

func Err() ProdPrinter {
	return printerErr
}

func Log() DevPrinter {
	return printerLog
}

func Debug() DevPrinter {
	return printerDebug
}

func DebugBatsTestBody() DevPrinter {
	return printerBatsTestBody
}

func DebugAllowCommit() DevPrinter {
	return printerDebug
}
