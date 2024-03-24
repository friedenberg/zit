package store_objekten

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/erworben"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/kilo/objekte_store"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/mike/store_util"
)

type Store struct {
	store_util.StoreUtil

	protoZettel      zettel.ProtoZettel
	konfigAkteFormat objekte.AkteFormat[erworben.Akte, *erworben.Akte]

	objekte_store.LogWriter
}

func Make(
	su store_util.StoreUtil,
) (s *Store, err error) {
	s = &Store{
		StoreUtil: su,
	}

	su.SetMatchableAdder(s)

	s.protoZettel = zettel.MakeProtoZettel(su.GetKonfig())

	s.konfigAkteFormat = objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
		objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](
			s.GetStandort(),
		),
		objekte.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
		s.GetStandort(),
	)

	errors.TodoP1("implement for other gattung")

	return
}

func (s *Store) SetLogWriter(lw objekte_store.LogWriter) {
	s.LogWriter = lw
}

func (s *Store) GetKonfigAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte] {
	return s.konfigAkteFormat
}

func (s Store) Flush(
	printerHeader schnittstellen.FuncIter[string],
) (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	if err = s.StoreUtil.Flush(printerHeader); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetAbbrStore().Flush(); err != nil {
		errors.Err().Print(err)
		err = errors.Wrapf(err, "failed to flush abbr index")
		return
	}

	return
}

func (s *Store) ReindexOne(besty, sk *sku.Transacted) (err error) {
	errExists := s.StoreUtil.GetAbbrStore().Exists(&sk.Kennung)

	if err = s.NewOrUpdated(errExists)(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.handleNewOrUpdatedCommit(
		sk,
		objekte_mode.ModeEmpty,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.AddTypToIndex(&sk.Metadatei.Typ); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetAbbrStore().AddMatchable(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P2 add support for quiet reindexing
func (s *Store) Reindex() (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = s.ResetIndexes(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetStandort().ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetVerzeichnisse().Initialize(
		s.GetKennungIndex(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(
		s.ReindexOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
