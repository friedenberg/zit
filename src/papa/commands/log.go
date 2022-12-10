package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Log struct {
}

func init() {
	registerCommand(
		"log",
		func(f *flag.FlagSet) Command {
			c := &Log{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

// TODO modify this to support other identifiers and provide option to search all
// or just schwanzen
func (c Log) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
	)

	return
}

func (c Log) RunWithIds(os *umwelt.Umwelt, is id_set.Set) (err error) {
	hs := is.Hinweisen()

	switch hs.Len() {

	case 0:
		err = errors.Errorf("hinweis or zettel sha required")
		return
	}

	//TODO-P2 switch to streams
	chains := make([][]*zettel.Transacted, 0, hs.Len())

	for _, h := range hs.Elements() {
		var chain []*zettel.Transacted

		if chain, err = os.StoreObjekten().Zettel().AllInChain(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		chains = append(chains, chain)
	}

	var b []byte

	if b, err = json.Marshal(chains); err != nil {
		err = errors.Wrapf(err, "failed to marshal json")
		return
	}

	errors.Out().Print(string(b))

	return
}
