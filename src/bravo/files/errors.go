package files

type ErrEmptyFileList struct{}

func (e ErrEmptyFileList) Error() string {
	return "empty file list"
}

func (e ErrEmptyFileList) Is(target error) (ok bool) {
	_, ok = target.(ErrEmptyFileList)
	return
}
