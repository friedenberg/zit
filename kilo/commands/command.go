package commands

import (
	"flag"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/konfig"
	"github.com/friedenberg/zit/delta/umwelt"
)

type Command interface {
	Run(*umwelt.Umwelt, ...string) error
}

type CommandSupportingErrors interface {
	HandleError(*umwelt.Umwelt, error)
}

type CommandWithArgPreprocessor interface {
	PreprocessArgs(*umwelt.Umwelt, []string) ([]string, error)
}

type CommandWithDescription interface {
	Description() string
}

type command struct {
	Command
	*flag.FlagSet
	DirZit string
}

var (
	commands = map[string]command{}
)

func Commands() map[string]command {
	return commands
}

func registerCommand(n string, makeFunc func(*flag.FlagSet) Command) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c := makeFunc(f)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = command{
		Command: c,
		FlagSet: f,
	}

	return
}

func Run(args []string) (exitStatus int) {
	var err error

	defer stdprinter.WaitForPrinter()
	defer func() {
		l := open_file_guard.Len()

		if l > 0 {
			stdprinter.Errf("file guard channel: %d\n", l)
		}

		var normalError errors.StackTracer

		if err != nil {
			//TODO use error to generate more specific exit status
			exitStatus = 1
		}

		if errors.As(err, &normalError) {
			stdprinter.Errf("%s\n", normalError.Error())
		} else {
			if err != nil {
				stdprinter.Error(err)
			}
		}
	}()

	var cmd command

	if err != nil {
		err = errors.Error(err)
		return
	}

	if len(os.Args) < 1 {
		logz.Print("printing usage")
		return cmd.PrintUsage(nil)
	}

	if len(os.Args) == 1 {
		return cmd.PrintUsage(errors.Errorf("No subcommand provided."))
	}

	cmds := Commands()
	specifiedSubcommand := os.Args[1]

	ok := false

	if cmd, ok = cmds[specifiedSubcommand]; !ok {
		return cmd.PrintUsage(errors.Errorf("No subcommand '%s'", specifiedSubcommand))
	}

	args = os.Args[2:]

	konfigCli := konfig.DefaultCli()
	konfigCli.AddToFlags(cmd.FlagSet)

	if err = cmd.FlagSet.Parse(args); err != nil {
		err = errors.Error(err)
		return
	}

	if konfigCli.Debug {
		df := cmd.SetDebug()
		defer df()
	}

	var k konfig.Konfig

	if k, err = konfigCli.Konfig(); err != nil {
		err = errors.Error(err)
		return
	}

	var u *umwelt.Umwelt

	if u, err = umwelt.MakeUmwelt(k); err != nil {
		err = errors.Error(err)
		return
	}

	cmdArgs := cmd.FlagSet.Args()

	if t, ok := cmd.Command.(CommandWithArgPreprocessor); ok {
		if cmdArgs, err = t.PreprocessArgs(u, cmdArgs); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if err = cmd.Command.Run(u, cmdArgs...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (c command) SetDebug() (d func()) {
	df := make([]func(), 0)
	d = func() {
		for i := len(df) - 1; i >= 0; i-- {
			df[i]()
		}
	}

	debug.SetGCPercent(-1)

	f, _ := open_file_guard.Create("build/cpu1.pprof")
	df = append(df, func() { open_file_guard.Close(f) })

	f1, _ := open_file_guard.Create("build/trace")
	df = append(df, func() { open_file_guard.Close(f1) })

	pprof.StartCPUProfile(f)
	df = append(df, func() { pprof.StopCPUProfile() })

	trace.Start(f1)
	df = append(df, func() { trace.Stop() })

	return
}
