package main

import (
	"os"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/kilo/commands"
)

func main() {
	exitStatus := commands.Run(os.Args)
	os.Exit(exitStatus)
	logz.Print()
}
