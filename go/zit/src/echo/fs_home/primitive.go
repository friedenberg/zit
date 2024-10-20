package fs_home

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
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
	sv       interfaces.StoreVersion
}

func MakePrimitive(do debug.Options) (s Primitive, err error) {
	if s.cwd, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.pid = os.Getpid()
	s.dryRun = do.DryRun

	if err = s.xdg.Initialize(true, "zit"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.sv, err = immutable_config.ReadStoreVersionFromFile(
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

func (h Primitive) DataFileStoreVersion() string {
	return filepath.Join(h.xdg.Data, "version")
}
