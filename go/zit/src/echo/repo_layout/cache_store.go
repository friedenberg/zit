package repo_layout

import "code.linenisgreat.com/zit/go/zit/src/delta/sha"

func (s Layout) ReadCloserCache(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             s.age,
		Path:            p,
		CompressionType: s.config.compressionType,
	}

	return NewFileReader(o)
}

func (s Layout) WriteCloserCache(
	p string,
) (w sha.WriteCloser, err error) {
	return NewMover(
		MoveOptions{
			Age:             s.age,
			FinalPath:       p,
			LockFile:        false,
			CompressionType: s.config.compressionType,
			TemporaryFS:     s.TempLocal,
		},
	)
}
