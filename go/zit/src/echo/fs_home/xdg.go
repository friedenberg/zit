package fs_home

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type XDG struct {
	Home      string
	AddedPath string
	Data      string
	Config    string
	State     string
	Cache     string
	Runtime   string
}

type xdgInitElement struct {
	defawlt string
	envKey  string
	out     *string
}

func (x *XDG) setDefaultOrEnv(
	initElement xdgInitElement,
) (err error) {
	if v, ok := os.LookupEnv(initElement.envKey); ok {
		*initElement.out = v
	} else {
		*initElement.out = os.Expand(initElement.defawlt, func(v string) string {
      switch v {
      case "HOME":
        return x.Home

      default:
        return os.Getenv(v)
      }
    })
	}

	*initElement.out = filepath.Join(*initElement.out, x.AddedPath)

	return
}

func (x *XDG) Initialize(mkDir bool, addedPath string) (err error) {
	x.AddedPath = addedPath

	toInitialize := []xdgInitElement{
		{
			defawlt: "$HOME/.local/share",
			envKey:  "XDG_DATA_HOME",
			out:     &x.Data,
		},
		{
			defawlt: "$HOME/.config",
			envKey:  "XDG_CONFIG_HOME",
			out:     &x.Config,
		},
		{
			defawlt: "$HOME/.local/state",
			envKey:  "XDG_STATE_HOME",
			out:     &x.State,
		},
		{
			defawlt: "$HOME/.cache",
			envKey:  "XDG_CACHE_HOME",
			out:     &x.Cache,
		},
		{
			defawlt: "$HOME/.local/runtime",
			envKey:  "XDG_RUNTIME_HOME",
			out:     &x.Cache,
		},
	}

	for _, ie := range toInitialize {
		if err = x.setDefaultOrEnv(ie); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = os.MkdirAll(*ie.out, 0o700); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
