package id

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/zk_types"
)

type Id interface {
	Kopf() string
	Schwanz() string
	String() string
	Sha() sha.Sha
}

type MutableId interface {
	Id
	flag.Value
}

type TypedId interface {
	Id
	Type() zk_types.Type
}

func Path(i Id, pc ...string) string {
	pc = append(pc, i.Kopf(), i.Schwanz())
	return path.Join(pc...)
}

func MakeDirIfNecessary(i Id, pc ...string) (p string, err error) {
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
