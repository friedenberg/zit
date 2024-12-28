package dir_layout

import (
	"fmt"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
)

const (
	EnvDir = "DIR_ZIT"
	EnvBin = "BIN_ZIT"
)

type Layout struct {
	*errors.Context
	beforeXDG
	xdg.XDG
}

func MakeDefault(
	context *errors.Context,
	do debug.Options,
) Layout {
	var home string

	{
		var err error

		if home, err = os.UserHomeDir(); err != nil {
			context.CancelWithError(err)
		}
	}

	return MakeWithHome(context, home, do, true)
}

func MakeDefaultAndInitialize(
	context *errors.Context,
	do debug.Options,
	overrideXDG bool,
) Layout {
	var home string

	{
		var err error
		if home, err = os.UserHomeDir(); err != nil {
			context.CancelWithError(err)
		}
	}

	return MakeWithHomeAndInitialize(
		context,
		home,
		do,
		overrideXDG,
	)
}

func MakeWithHome(
	context *errors.Context,
	home string,
	do debug.Options,
	permitCwdXDGOverride bool,
) (s Layout) {
	s.Context = context

	xdg := xdg.XDG{
		Home: home,
	}

	if err := s.beforeXDG.initialize(do); err != nil {
		s.CancelWithError(err)
	}

	addedPath := "zit"
	pathCwdXDGOverride := filepath.Join(s.cwd, ".zit")

	if permitCwdXDGOverride && files.Exists(pathCwdXDGOverride) {
		xdg.Home = pathCwdXDGOverride
		addedPath = ""
		if err := xdg.InitializeOverridden(false, addedPath); err != nil {
			s.CancelWithError(err)
		}
	} else {
		if err := xdg.InitializeStandardFromEnv(false, addedPath); err != nil {
			s.CancelWithError(err)
		}
	}

	if err := s.initializeXDG(xdg); err != nil {
		s.CancelWithError(err)
	}

	s.AfterWithContext(s.resetTempOnExit)

	return
}

func MakeWithHomeAndInitialize(
	context *errors.Context,
	home string,
	do debug.Options,
	cwdXDGOverride bool,
) (s Layout) {
	s.Context = context

	xdg := xdg.XDG{
		Home: home,
	}

	if err := s.beforeXDG.initialize(do); err != nil {
		s.CancelWithError(err)
	}

	addedPath := "zit"
	pathCwdXDGOverride := filepath.Join(s.cwd, ".zit")

	if cwdXDGOverride {
		xdg.Home = pathCwdXDGOverride
		addedPath = ""
		if err := xdg.InitializeOverridden(true, addedPath); err != nil {
			s.CancelWithError(err)
		}
	} else {
		if err := xdg.InitializeStandardFromEnv(true, addedPath); err != nil {
			s.CancelWithError(err)
		}
	}

	if err := s.initializeXDG(xdg); err != nil {
		s.CancelWithError(err)
	}

	s.AfterWithContext(s.resetTempOnExit)

	return
}

func MakeWithXDG(
	context *errors.Context,
	do debug.Options,
	xdg xdg.XDG,
) (s Layout) {
	s.Context = context

  if err := s.beforeXDG.initialize(do); err != nil {
    s.CancelWithError(err)
	}

  if err := s.initializeXDG(xdg); err != nil {
    s.CancelWithError(err)
	}

	return
}

func (layout *Layout) initializeXDG(xdg xdg.XDG) (err error) {
	layout.XDG = xdg

	layout.TempLocal.BasePath = filepath.Join(
		layout.Cache,
		fmt.Sprintf("tmp-%d", layout.GetPid()),
	)

	if err = layout.MakeDir(layout.TempLocal.BasePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h Layout) GetDebug() debug.Options {
	return h.debug
}

func (h Layout) IsDryRun() bool {
	return h.dryRun
}

func (h Layout) GetPid() int {
	return h.pid
}

func (h Layout) GetExecPath() string {
	return h.execPath
}

func (h Layout) GetCwd() string {
	return h.cwd
}

func (h Layout) GetXDG() xdg.XDG {
	return h.XDG
}

func (h *Layout) SetXDG(x xdg.XDG) {
	h.XDG = x
}

func (s Layout) AbsFromCwdOrSame(p string) (p1 string) {
	var err error
	p1, err = filepath.Abs(p)
	if err != nil {
		p1 = p
	}

	return
}

func (s Layout) RelToCwdOrSame(p string) (p1 string) {
	var err error

	if p1, err = filepath.Rel(s.GetCwd(), p); err != nil {
		p1 = p
	}

	return
}

func (s Layout) Rel(
	p string,
) (out string) {
	out = p

	p1, _ := filepath.Rel(s.GetCwd(), p)

	if p1 != "" {
		out = p1
	}

	return
}

func (h Layout) MakeCommonEnv() map[string]string {
	return map[string]string{
		"ZIT_BIN": h.GetExecPath(),
		// TODO determine if ZIT_DIR is kept
		// "ZIT_DIR": h.Dir(),
	}
}

func (s Layout) MakeDir(ds ...string) (err error) {
	for _, d := range ds {
		if err = os.MkdirAll(d, os.ModeDir|0o755); err != nil {
			err = errors.Wrapf(err, "Dir: %q", d)
			return
		}
	}

	return
}
