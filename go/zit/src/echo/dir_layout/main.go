package dir_layout

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

const (
	EnvDir = "DIR_ZIT"
	EnvBin = "BIN_ZIT"
)

type Layout struct {
	cwd      string
	execPath string
	dryRun   bool
	debug    debug.Options
	pid      int
	xdg      xdg.XDG
	sv       immutable_config.StoreVersion
}

func MakePrimitive(do debug.Options) (s Layout, err error) {
	var home string

	if home, err = os.UserHomeDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s, err = MakePrimitiveWithHome(home, do); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakePrimitiveWithHome(
	home string,
	do debug.Options,
) (s Layout, err error) {
	xdg := xdg.XDG{
		Home: home,
	}

	if err = xdg.InitializeFromEnv(true, "zit"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s, err = MakePrimitiveWithXDG(do, xdg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakePrimitiveWithXDG(
	do debug.Options,
	xdg xdg.XDG,
) (s Layout, err error) {
	if s.cwd, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.pid = os.Getpid()
	s.dryRun = do.DryRun

	s.xdg = xdg

	if err = s.sv.ReadFromFile(
		s.DataFileStoreVersion(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.execPath, err = os.Executable(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO switch to useing MakeCommonEnv()
	{
		if err = os.Setenv(EnvBin, s.execPath); err != nil {
			err = errors.Wrap(err)
			return
		}
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

func (h Layout) GetStoreVersion() immutable_config.StoreVersion {
	return h.sv
}

func (h Layout) DataFileStoreVersion() string {
	return filepath.Join(h.xdg.Data, "version")
}

func (h Layout) GetXDG() xdg.XDG {
	return h.xdg
}

func (h *Layout) SetXDG(x xdg.XDG) {
	h.xdg = x
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
