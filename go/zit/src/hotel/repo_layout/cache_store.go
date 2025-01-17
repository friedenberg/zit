package repo_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
)

func (s Layout) ReadCloserCache(p string) (sha.ReadCloser, error) {
	o := env_dir.FileReadOptions{
		// Config: s.Config.Blob,
		Path: p,
	}

	return env_dir.NewFileReader(o)
}

func (s Layout) WriteCloserCache(
	p string,
) (w sha.WriteCloser, err error) {
	return env_dir.NewMover(
		env_dir.MoveOptions{
			// Config:      s.Config.Blob,
			FinalPath:   p,
			TemporaryFS: s.GetTempLocal(),
		},
	)
}
