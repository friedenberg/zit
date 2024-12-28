package repo_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

func (s Layout) Initialize() {
	if err := s.MakeDir(
		s.DirObjectId(),
		s.DirCache(),
		s.DirLostAndFound(),
	); err != nil {
		s.CancelWithError(err)
	}

	for _, g := range []genres.Genre{genres.Blob, genres.InventoryList} {
		var d string
		var err error

		if d, err = s.DirObjectGenre(g); err != nil {
			if genres.IsErrUnsupportedGenre(err) {
				err = nil
				continue
			} else {
				s.CancelWithError(err)
			}
		}

		if err := s.MakeDir(d); err != nil {
			s.CancelWithError(err)
		}
	}
}
