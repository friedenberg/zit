package main

import (
	"os"
	"os/signal"
	"syscall"

	"code.linenisgreat.com/zit-go/src/quebec/commands"
)

func main() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)

	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c, syscall.SIGURG)

	go func() {
		for range c {
			// logz.Printf("signal: %s", s)
		}
	}()

	exitStatus := commands.Run(os.Args)
	os.Exit(exitStatus)
}
