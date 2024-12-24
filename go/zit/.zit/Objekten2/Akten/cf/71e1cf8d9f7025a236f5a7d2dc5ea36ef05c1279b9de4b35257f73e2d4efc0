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
	defawlt string
	envKey  string
	out     *string
}

func (x *XDG) GetInitElements() []xdgInitElement {
	return []xdgInitElement{
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
			out:     &x.Runtime,
		},
	}
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

func (x *XDG) InitializeFromEnv(mkDir bool, addedPath string) (err error) {
	x.AddedPath = addedPath

	toInitialize := x.GetInitElements()

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

func (x *XDG) InitializeFromFile(
	mkDir bool,
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
