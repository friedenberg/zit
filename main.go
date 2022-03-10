package main

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/juliett/commands"
)

func main() {
	if err := commands.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
