package repo_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

func (s Layout) ReadCloserCache(p string) (sha.ReadCloser, error) {
	o := dir_layout.FileReadOptions{
		// Config: s.Config.Blob,
		Path: p,
	}

	return dir_layout.NewFileReader(o)
}

func (s Layout) WriteCloserCache(
	p string,
) (w sha.WriteCloser, err error) {
	return dir_layout.NewMover(
		dir_layout.MoveOptions{
			// Config:      s.Config.Blob,
			FinalPath:   p,
			TemporaryFS: s.TempLocal,
		},
	)
}
