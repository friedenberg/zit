package lua

import (
	"io"
	"os"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/files"
)

type LuaFlag struct {
	VMPool
	value string
}

func (l *LuaFlag) Set(script string) (err error) {
	l.value = script

	var f *os.File

	f, err = files.Open(script)

	if errors.IsNotExist(err) {
		err = nil
	} else if !errors.IsNotExist(err) && err != nil {
		err = errors.Wrap(err)
		return
	} else if err == nil {
		var sb strings.Builder

		if _, err = io.Copy(&sb, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		script = sb.String()

		if err = f.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = l.VMPool.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (lf *LuaFlag) String() string {
	return lf.value
}
