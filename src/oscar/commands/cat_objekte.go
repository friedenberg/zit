package commands

import (
	"flag"
	"fmt"
	"io"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	gattung "github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/india/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/mike/umwelt"
)

type CatObjekte struct {
	Type gattung.Gattung
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{
				Type: gattung.Unknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c CatObjekte) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &konfig.Id{},
		},
		id_set.ProtoId{
			MutableId: &sha.Sha{},
		},
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &etikett.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e etikett.Etikett
				e, err = u.StoreObjekten().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &typ.Typ{},
		},
		id_set.ProtoId{
			MutableId: &ts.Time{},
		},
	)

	return
}

// TODO switch to idset semantics
func (c CatObjekte) RunWithIds(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	switch c.Type {

	case gattung.Akte:
		return c.akten(store, ids)

	case gattung.Zettel:
		return c.zettelen(store, ids)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c CatObjekte) akteShasFromIds(
	u *umwelt.Umwelt,
	ids id_set.Set,
) (zettelen zettel_transacted.MutableSet, err error) {
	zettelen = zettel_transacted.MakeMutableSetUnique(0)

	if err = u.StoreObjekten().ReadAllSchwanzenVerzeichnisse(
		zettel_verzeichnisse.WriterZettelTransacted{
			Writer: zettel_transacted.WriterZettelNamed{
				Writer: zettel_named.FilterIdSet{
					Set: ids,
				},
			},
		},
		zettel_verzeichnisse.WriterZettelTransacted{
			Writer: zettel_transacted.MakeWriter(zettelen.Add),
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// for _, h := range ids.Hinweisen().Elements() {
	// 	var zc zettel_checked_out.Zettel

	// 	if zc, err = u.StoreWorkingDirectory().Read(h.String() + ".md"); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	if zc.State == zettel_checked_out.StateExistsAndDifferent {
	// 		shas = append(shas, zc.External.Named.Stored.Zettel.Akte)
	// 	} else {
	// 		shas = append(shas, zc.Internal.Named.Stored.Zettel.Akte)
	// 	}
	// }

	return
}

func (c CatObjekte) akten(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	var zettelen zettel_transacted.MutableSet

	if zettelen, err = c.akteShasFromIds(store, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	zettelen.Each(
		func(z *zettel_transacted.Zettel) (err error) {
			var r io.ReadCloser

			sb := z.Named.Stored.Zettel.Akte

			if r, err = store.StoreObjekten().AkteReader(sb); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.PanicIfError(r.Close)

			if _, err = io.WriteString(
				store.Out(),
				fmt.Sprintf("['%s']\n", z.Named.Hinweis),
			); err != nil {
				err = errors.IsAsNilOrWrapf(
					err,
					syscall.EPIPE,
					"Zettel: %s",
					z.Named.Hinweis,
				)

				return
			}

			if _, err = io.Copy(store.Out(), r); err != nil {
				err = errors.IsAsNilOrWrapf(
					err,
					syscall.EPIPE,
					"Zettel: %s",
					z.Named.Hinweis,
				)

				return
			}

			return
		},
	)

	return
}

func (c CatObjekte) zettelen(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	w := zettel_transacted.MakeWriterChain(
		zettel_transacted.WriterZettelNamed{
			Writer: zettel_named.WriterFilter{
				NamedFilter: zettel_named.FilterIdSet{
					Set: ids,
				},
			},
		},
		zettel_transacted.MakeWriterZettel(
			zettel.MakeSerializedFormatWriter(
				&zettel.Objekte{},
				store.Out(),
				store.StoreObjekten(),
				store.Konfig(),
			),
		),
	)

	if err = store.StoreObjekten().ReadAllSchwanzenTransacted(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
