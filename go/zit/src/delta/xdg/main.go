package xdg

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

type XDG struct {
	Home      string
	AddedPath string // name of the utility

	Data    string
	Config  string
	State   string
	Cache   string
	Runtime string
}

type xdgInitElement struct {
	standard   string
	overridden string
	envKey     string
	out        *string
}

func (x *XDG) GetInitElements() []xdgInitElement {
	return []xdgInitElement{
		{
			standard:   "$HOME/.local/share",
			overridden: "$HOME/local/share",
			envKey:     "XDG_DATA_HOME",
			out:        &x.Data,
		},
		{
			standard:   "$HOME/.config",
			overridden: "$HOME/config",
			envKey:     "XDG_CONFIG_HOME",
			out:        &x.Config,
		},
		{
			standard:   "$HOME/.local/state",
			overridden: "$HOME/local/state",
			envKey:     "XDG_STATE_HOME",
			out:        &x.State,
		},
		{
			standard:   "$HOME/.cache",
			overridden: "$HOME/cache",
			envKey:     "XDG_CACHE_HOME",
			out:        &x.Cache,
		},
		{
			standard:   "$HOME/.local/runtime",
			overridden: "$HOME/local/runtime",
			envKey:     "XDG_RUNTIME_HOME",
			out:        &x.Runtime,
		},
	}
}

func (x *XDG) setDefaultOrEnv(
	defaultValue string,
	envKey string,
) (out string, err error) {
	if v, ok := os.LookupEnv(envKey); envKey != "" && ok {
		out = v
	} else {
		out = os.Expand(
			defaultValue,
			func(v string) string {
				switch v {
				case "HOME":
					return x.Home

				default:
					return os.Getenv(v)
				}
			},
		)
	}

	if x.AddedPath != "" {
		out = filepath.Join(out, x.AddedPath)
	}

	return
}

func (x *XDG) InitializeOverridden(
	mkDir bool,
	addedPath string,
) (err error) {
	x.AddedPath = addedPath

	toInitialize := x.GetInitElements()

	for _, ie := range toInitialize {
		if *ie.out, err = x.setDefaultOrEnv(
			ie.overridden,
			"",
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if mkDir {
			if err = os.MkdirAll(*ie.out, 0o700); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (x *XDG) InitializeStandardFromEnv(
	mkDir bool,
	addedPath string,
) (err error) {
	x.AddedPath = addedPath

	toInitialize := x.GetInitElements()

	for _, ie := range toInitialize {
		if *ie.out, err = x.setDefaultOrEnv(
			ie.standard,
			ie.envKey,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if mkDir {
			if err = os.MkdirAll(*ie.out, 0o700); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (x *XDG) InitializeFromDotenvFile(
	addedPath string,
	file string,
) (err error) {
	x.AddedPath = addedPath

	var f *os.File

	if f, err = files.Open(file); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	r := Dotenv{
		XDG: x,
	}

	if _, err = r.ReadFrom(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
