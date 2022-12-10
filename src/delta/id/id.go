package id

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/sha"
)

type Id interface {
	String() string
	Sha() sha.Sha
}

type IdMitKorper interface {
	Id
	Kopf() string
	Schwanz() string
}

type MutableId interface {
	Id
	flag.Value
}

type TypedId interface {
	Id
	Type() gattung.Gattung
}

func Path(i IdMitKorper, pc ...string) string {
	pc = append(pc, i.Kopf(), i.Schwanz())
	return path.Join(pc...)
}

func MakeDirIfNecessary(i IdMitKorper, pc ...string) (p string, err error) {
	p = Path(i, pc...)
	dir := path.Dir(p)

	if err = os.MkdirAll(dir, os.ModeDir|0755); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

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
