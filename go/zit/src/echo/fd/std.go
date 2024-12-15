package fd

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

type Std struct {
	*os.File
	IsTty bool
}

func MakeStd(f *os.File) Std {
	return Std{
		File:  f,
		IsTty: files.IsTty(f),
	}
}

func (f Std) GetFile() *os.File {
	return f.File
}
