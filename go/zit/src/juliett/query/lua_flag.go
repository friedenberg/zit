package query

import (
	"io"
	"os"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/charlie/files"
	lua "github.com/yuin/gopher-lua"
	lua_ast "github.com/yuin/gopher-lua/ast"
	lua_parse "github.com/yuin/gopher-lua/parse"
)

type LuaFlag struct {
	Lua
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

	reader := strings.NewReader(script)

	var chunks []lua_ast.Stmt

	if chunks, err = lua_parse.Parse(reader, ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	var compiled *lua.FunctionProto

	if compiled, err = lua.Compile(chunks, ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	l.statePool = pool.MakePool(
		func() (l *luaSku) {
			l = &luaSku{
				LState: lua.NewState(),
			}

			l.LTable = l.NewTable()

			lfunc := l.NewFunctionFromProto(compiled)
			l.Push(lfunc)
			l.PCall(0, lua.MultRet, nil)

			return l
		},
		func(s *luaSku) {
			s.SetTop(0)
		},
	)

	return
}

func (lf *LuaFlag) String() string {
	return lf.value
}
