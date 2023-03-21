package store_fs

import (
	"fmt"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
)

func (s *Store) FileExtensionForGattung(
	gg schnittstellen.GattungGetter,
) string {
	switch gattung.Must(gg.GetGattung()) {
	case gattung.Zettel:
		return s.erworben.FileExtensions.Zettel

	case gattung.Etikett:
		return s.erworben.FileExtensions.Etikett

	case gattung.Typ:
		return s.erworben.FileExtensions.Etikett

	default:
		return "unknown_gattung"
	}
}

func (s *Store) PathForTransactedLike(tl objekte.TransactedLike) string {
	return path.Join(
		s.Cwd(),
		fmt.Sprintf(
			"%s.%s",
			tl.GetDataIdentity().GetId(),
			s.FileExtensionForGattung(tl.GetDataIdentity()),
		),
	)
}

func (s *Store) Query(
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.CheckedOutLike],
) (err error) {
	if err = s.storeObjekten.Query(
		ms,
		func(t objekte.TransactedLike) (err error) {
			var co objekte.CheckedOutLike

			if co, err = s.readOneGeneric(t); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			return f(co)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readOneGeneric(
	t objekte.TransactedLike,
) (co objekte.CheckedOutLike, err error) {
	switch tt := t.(type) {
	case *zettel.Transacted:
		return s.ReadOneZettel(*tt)

	case *typ.Transacted:
		co, err = s.ReadOneTyp(*tt)

	case *etikett.Transacted:
		co, err = s.ReadOneEtikett(*tt)

	default:
		err = gattung.MakeErrUnsupportedGattung(tt.GetSku2())
		return
	}

	return
}

func (s *Store) ReadOneZettel(
	sz zettel.Transacted,
) (cz zettel.CheckedOut, err error) {
	p := s.PathForTransactedLike(sz)

	if cz, err = s.readOneFS(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	cz.Internal = sz

	return
}

func (s *Store) ReadOneEtikett(
	tk etikett.Transacted,
) (co etikett.CheckedOut, err error) {
	co.Internal = tk
	co.External.Sku = tk.Sku.GetExternal()

	p := s.PathForTransactedLike(tk)

	if co.External, err = s.storeObjekten.Etikett().ReadOneExternal(
		cwd.Etikett{
			Kennung: tk.Sku.Kennung,
			FDs: sku.ExternalFDs{
				Objekte: kennung.FD{
					Path: p,
				},
			},
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.External.Sku.Kennung = tk.Sku.Kennung

	return
}

func (s *Store) ReadOneTyp(
	tk typ.Transacted,
) (co typ.CheckedOut, err error) {
	co.Internal = tk
	co.External.Sku = tk.Sku.GetExternal()

	p := s.PathForTransactedLike(tk)

	if co.External, err = s.ReadTyp(
		cwd.Typ{
			Kennung: tk.Sku.Kennung,
			FDs: sku.ExternalFDs{
				Objekte: kennung.FD{
					Path: p,
				},
			},
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.External.Sku.Kennung = tk.Sku.Kennung

	return
}
