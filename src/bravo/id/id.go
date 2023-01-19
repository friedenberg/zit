package id

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type TypedId interface {
	schnittstellen.Gattung
	schnittstellen.Setter
}

// func Path(i IdMitKorper, pc ...string) string {
// 	pc = append(pc, i.Kopf(), i.Schwanz())
// 	return path.Join(pc...)
// }

func Path(i schnittstellen.Korper, pc ...string) string {
	pc = append(pc, i.Kopf(), i.Schwanz())
	return path.Join(pc...)
}

func MakeDirIfNecessary(i schnittstellen.Korper, pc ...string) (p string, err error) {
	p = Path(i, pc...)
	dir := path.Dir(p)

	if err = os.MkdirAll(dir, os.ModeDir|0755); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P2 determine if this is used
func HeadTailFromFileName(fileName string) (head string, tail string) {
	head, tail = filepath.Split(fileName)
	tail = tail[0 : len(tail)-len(path.Ext(tail))]

	head = strings.TrimSuffix(head, string(filepath.Separator))

	idx := strings.LastIndex(head, string(filepath.Separator))

	if idx != -1 {
		head = head[idx+1:]
	}

	return
}
