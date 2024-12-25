package dir_layout_primitive

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

type Primitive struct {
	cwd      string
	execPath string
	dryRun   bool
	debug    debug.Options
	pid      int
	xdg      xdg.XDG
	sv       immutable_config.StoreVersion
}

func MakePrimitive(do debug.Options) (s Primitive, err error) {
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
) (s Primitive, err error) {
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
) (s Primitive, err error) {
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

func (h Primitive) GetDebug() debug.Options {
	return h.debug
}

func (h Primitive) IsDryRun() bool {
	return h.dryRun
}

func (h Primitive) GetPid() int {
	return h.pid
}

func (h Primitive) GetExecPath() string {
	return h.execPath
}

func (h Primitive) GetCwd() string {
	return h.cwd
}

func (h Primitive) GetStoreVersion() immutable_config.StoreVersion {
	return h.sv
}

func (h Primitive) DataFileStoreVersion() string {
	return filepath.Join(h.xdg.Data, "version")
}

func (h Primitive) GetXDG() xdg.XDG {
	return h.xdg
}

func (h *Primitive) SetXDG(x xdg.XDG) {
	h.xdg = x
}

func (s Primitive) AbsFromCwdOrSame(p string) (p1 string) {
	var err error
	p1, err = filepath.Abs(p)
	if err != nil {
		p1 = p
	}

	return
}

func (s Primitive) RelToCwdOrSame(p string) (p1 string) {
	var err error

	if p1, err = filepath.Rel(s.GetCwd(), p); err != nil {
		p1 = p
	}

	return
}

func (s Primitive) Rel(
	p string,
) (out string) {
	out = p

	p1, _ := filepath.Rel(s.GetCwd(), p)

	if p1 != "" {
		out = p1
	}

	return
}
