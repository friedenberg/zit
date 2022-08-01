package main

import (
	"os"

	"github.com/friedenberg/zit/kilo/commands"
)

func main() {
	exitStatus := commands.Run(os.Args)
	os.Exit(exitStatus)
}
