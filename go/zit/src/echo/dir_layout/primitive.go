package dir_layout

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

type Primitive struct {
	cwd      string
	execPath string
	dryRun   bool
	debug    debug.Options
	pid      int
	xdg      XDG
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

func MakePrimitiveWithHome(home string, do debug.Options) (s Primitive, err error) {
	if s.cwd, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.pid = os.Getpid()
	s.dryRun = do.DryRun

	s.xdg.Home = home

	if err = s.xdg.Initialize(true, "zit"); err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (h Primitive) GetStoreVersion() immutable_config.StoreVersion {
	return h.sv
}

func (h Primitive) DataFileStoreVersion() string {
	return filepath.Join(h.xdg.Data, "version")
}
