package commands

var commands = map[string]Command{}

func Commands() map[string]Command {
	return commands
}

func registerCommand(n string, cmd Command) {
	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = cmd
}
