package commands

import "code.linenisgreat.com/zit/go/zit/src/golf/command"

var commands = map[string]command.Command{}

func Commands() map[string]command.Command {
	return commands
}

func registerCommand(n string, cmd command.Command) {
	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = cmd
}
