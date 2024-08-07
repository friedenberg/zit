package commands

import (
	"bufio"
	"bytes"
	"flag"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func (c command) PrintUsage(in error) (exitStatus int) {
	if in != nil {
		exitStatus = 1
		ui.Err().Print(in)
	}

	ui.Err().Print("Usage for z:")

	fs := make([]flag.FlagSet, 0, len(Commands()))

	for _, c := range Commands() {
		fs = append(fs, *c.FlagSet)
	}

	sort.Slice(fs, func(i, j int) bool {
		return fs[i].Name() < fs[j].Name()
	})

	for _, f := range fs {
		c.PrintSubcommandUsage(f)
	}

	return
}

func (c command) PrintSubcommandUsage(flags flag.FlagSet) {
	printTabbed := func(s string) {
		ui.Err().Printf("  %s", s)
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
