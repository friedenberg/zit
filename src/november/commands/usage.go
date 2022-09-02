package commands

import (
	"bufio"
	"bytes"
	"flag"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func (c command) PrintUsage(in error) (exitStatus int) {
	if in != nil {
		exitStatus = 1
		errors.PrintErr(in)
	}

	errors.PrintErr("Usage for z:")

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
		errors.PrintErrf("  %s\n", s)
	}

	var b bytes.Buffer
	flags.SetOutput(&b)

	printTabbed(flags.Name())

	//TODO determine why the interface doesn't actually work
	if cwd, ok := c.Command.(CommandWithDescription); ok {
		printTabbed(cwd.Description())
	}

	flags.PrintDefaults()

	scanner := bufio.NewScanner(&b)

	for scanner.Scan() {
		printTabbed(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		errors.PrintErr(err)
	}
}
