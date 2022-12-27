package commands

import (
	"encoding/gob"
	"flag"
	"net"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/remote_messages"
)

type Test1 struct {
}

func init() {
	registerCommand(
		"test1",
		func(f *flag.FlagSet) Command {
			c := &Test1{}

			return c
		},
	)
}

func (c Test1) Run(u *umwelt.Umwelt, args ...string) (err error) {
	sockPath := args[0]

	var conn net.Conn

	if conn, err = net.Dial("unix", sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, conn.Close)

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	msg := remote_messages.Message{
		MessageType: remote_messages.MessageTypeSenderHi,
	}

	for {
		if err = enc.Encode(msg); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = dec.Decode(&msg); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Err().Print(msg)

    if !msg.NextLine() {
      break
    }
	}

	return
}
