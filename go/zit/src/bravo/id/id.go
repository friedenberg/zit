package id

import (
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type TypedId interface {
	interfaces.GattungLike
	interfaces.Setter
}

// func Path(i IdMitKorper, pc ...string) string {
// 	pc = append(pc, i.Kopf(), i.Schwanz())
// 	return path.Join(pc...)
// }

func Path(i interfaces.Korper, pc ...string) string {
	pc = append(pc, i.Kopf(), i.Schwanz())
	return path.Join(pc...)
}

func MakeDirIfNecessary(i interfaces.Korper, pc ...string) (p string, err error) {
	p = Path(i, pc...)
	dir := path.Dir(p)

	if err = os.MkdirAll(dir, os.ModeDir|0o755); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
