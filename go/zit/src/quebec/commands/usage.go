package commands

import (
	"bufio"
	"bytes"
	"flag"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func PrintUsage(ctx *errors.Context, in error) {
	if in != nil {
		defer ctx.CancelWithError(in)
	}

	ui.Err().Print("Usage for zit:")

	fs := make([]*flag.FlagSet, 0, len(Commands()))

	for name, cmd := range Commands() {
		f := flag.NewFlagSet(name, flag.ExitOnError)
		cmd.SetFlagSet(f)
		fs = append(fs, f)
	}

	sort.Slice(fs, func(i, j int) bool {
		return fs[i].Name() < fs[j].Name()
	})

	for _, f := range fs {
		ui.Err().Print(f.Name())
	}
}

func PrintSubcommandUsage(flags flag.FlagSet) {
	printTabbed := func(s string) {
		ui.Err().Print(s)
	}

	var b bytes.Buffer
	flags.SetOutput(&b)

	printTabbed(flags.Name())

	flags.PrintDefaults()

	scanner := bufio.NewScanner(&b)

	for scanner.Scan() {
		printTabbed(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		ui.Err().Print(err)
	}
}
