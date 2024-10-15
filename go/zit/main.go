package main

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/quebec/commands"
)

func main() {
	exitStatus := commands.Run(os.Args)
	os.Exit(exitStatus)
}
