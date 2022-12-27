package commands

import (
	"encoding/gob"
	"flag"
	"net"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/remote_messages"
)

type Test struct {
}

func init() {
	registerCommand(
		"test",
		func(f *flag.FlagSet) Command {
			c := &Test{}

			return c
		},
	)
}

func (c Test) Run(u *umwelt.Umwelt, args ...string) (err error) {
	sockPath := args[0]

	var socket net.Listener

	if socket, err = net.Listen("unix", sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, func() error {
		return syscall.Unlink(sockPath)
	})

	// Accept an incoming connection.
	var conn net.Conn

	if conn, err = socket.Accept(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	msg := remote_messages.Message{
		MessageType: remote_messages.MessageTypeReceiverHi,
	}

	for {
		if err = dec.Decode(&msg); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		errors.Err().Print(msg)

		if !msg.NextLine() {
			break
		}

		if err = enc.Encode(msg); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
