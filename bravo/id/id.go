package id

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Id interface {
	Head() string
	Tail() string
	String() string
}

type MutableId interface {
	Id
	SetParts(string, string)
}

func Path(i Id, pc ...string) string {
	pc = append(pc, i.Head(), i.Tail())
	return path.Join(pc...)
}

func MakeDirIfNecessary(i Id, pc ...string) (p string, err error) {
	p = Path(i, pc...)
	dir := path.Dir(p)

	//TODO open_file_guard
	if err = os.MkdirAll(dir, os.ModeDir|0755); err != nil {
		err = _Error(err)
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
