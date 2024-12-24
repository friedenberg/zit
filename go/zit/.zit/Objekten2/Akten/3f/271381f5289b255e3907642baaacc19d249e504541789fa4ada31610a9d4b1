package dir_layout

import "code.linenisgreat.com/zit/go/zit/src/delta/sha"

func (s DirLayout) ReadCloserCache(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             s.age,
		Path:            p,
		CompressionType: s.immutable_config.CompressionType,
	}

	return NewFileReader(o)
}

func (s DirLayout) WriteCloserCache(
	p string,
) (w sha.WriteCloser, err error) {
	return NewMover(
		MoveOptions{
			Age:             s.age,
			FinalPath:       p,
			LockFile:        false,
			CompressionType: s.immutable_config.CompressionType,
			TemporaryFS:     s.TempLocal,
		},
	)
}
