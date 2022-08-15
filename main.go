package main

import (
	"os"
	"os/signal"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/kilo/commands"
)

func main() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)

	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c)

	go func() {
		for s := range c {
			logz.Printf("signal: %s", s)
		}
	}()

	exitStatus := commands.Run(os.Args)
	os.Exit(exitStatus)
	logz.Print()
}
