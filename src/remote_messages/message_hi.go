package remote_messages

import (
	"filippo.io/age"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type MessageHi struct {
	AgeRecipient string
	CliKonfig    konfig.Cli
}

func PerformHi(s *Stage, u *umwelt.Umwelt) (theirCliKonfig konfig.Cli, err error) {
	var i *age.X25519Identity

	{
		if i, err = age.GenerateX25519Identity(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	msgOurHi := MessageHi{
		AgeRecipient: i.Recipient().String(),
		CliKonfig:    u.Konfig().Cli(),
	}

	{
		if err = s.Send(msgOurHi); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var msgTheirHi MessageHi

	{
		if err = s.Receive(&msgTheirHi); err != nil {
			err = errors.Wrap(err)
			return
		}

		theirCliKonfig = msgTheirHi.CliKonfig
	}

	// var r age.Recipient

	{
		if _, err = age.ParseX25519Recipient(msgTheirHi.AgeRecipient); err != nil {
			err = errors.Wrap(err)
			return
		}

		// var w io.Writer

		// if w, err = age.Encrypt(s.conn, r); err != nil {
		// 	err = errors.Wrap(err)
		// 	return
		// }

		// enc = gob.NewEncoder(w)

		if err = s.Send(struct{}{}); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	{
		// var r1 io.Reader

		// if r1, err = age.Decrypt(s.conn, i); err != nil {
		// 	err = errors.Wrap(err)
		// 	return
		// }

		// errors.Err().Printf("did create reader")

		// s.dec = gob.NewDecoder(r1)

		// errors.Err().Print("will decode hi-ack")

		var msgTheirHi struct{}

		if err = s.Receive(&msgTheirHi); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
