package commands

import (
	"encoding/json"
	"flag"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type ReadAkte struct{}

func init() {
	registerCommand(
		"read-akte",
		func(f *flag.FlagSet) Command {
			c := &ReadAkte{}

			return c
		},
	)
}

type readAkteEntry struct {
	Akte string `json:"akte"`
}

func (c ReadAkte) Run(u *umwelt.Umwelt, args ...string) (err error) {
	dec := json.NewDecoder(u.In())

	for {
		var entry readAkteEntry

		if err = dec.Decode(&entry); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		var sh *sha.Sha

		if sh, err = c.readOneAkte(u, entry); err != nil {
			err = errors.Wrap(err)
			return
		}

		log.Debug().Print(sh)
	}

	return
}

func (ReadAkte) readOneAkte(u *umwelt.Umwelt, entry readAkteEntry) (sh *sha.Sha, err error) {
	var aw sha.WriteCloser

	if aw, err = u.Standort().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = io.Copy(aw, strings.NewReader(entry.Akte)); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.GetPool().Get()

	if err = sh.SetShaLike(aw.GetShaLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
