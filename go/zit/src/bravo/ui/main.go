package ui

import (
	"log"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

var (
	SetOutput = log.SetOutput
	verbose   bool
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
	errors.SetTesting()
	SetVerbose(true)
}

func IsVerbose() bool {
	return verbose
}

type Printer interface {
	GetPrinter() Printer
	GetFile() *os.File
	IsTty() bool
	Print(v ...interface{}) error
	Printf(format string, v ...interface{}) error
}

type DevPrinter interface {
	Printer
	Caller(i int, vs ...interface{})
	FunctionName(skip int)
}

var (
	printerOut, printerErr   printer
	printerLog, printerDebug devPrinter
	printerBatsTestBody      devPrinter
)

func init() {
	printerOut = MakePrinterOn(os.Stdout, true)
	printerErr = MakePrinterOn(os.Stderr, true)

	printerLog = devPrinter{
		printer:       printerErr.withOn(false),
		includesStack: true,
	}

	// TODO-P2 determine if on thru compilation
	printerDebug = devPrinter{
		printer:       printerErr,
		includesStack: true,
	}

	// TODO-P2 determine thru compilation
	envVarFilter := "BATS_TEST_BODY"
	_, printerBatsTestBodyOn := os.LookupEnv(envVarFilter)

	// TODO-P2 determine thru compilation
	printerBatsTestBody = devPrinter{
		printer:       printerErr.withOn(printerBatsTestBodyOn),
		includesStack: true,
	}
}

func Out() Printer {
	return printerOut
}

func Err() Printer {
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
