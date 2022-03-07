package commands

import (
	"flag"
)

type Edit struct {
	IncludeAkte bool
}

func init() {
	registerCommand(
		"edit",
		func(f *flag.FlagSet) Command {
			c := &Edit{}

			f.BoolVar(&c.IncludeAkte, "include-akte", true, "check out and open the akte")

			return commandWithZettels{c}
		},
	)
}

func (c Edit) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	var czs []_ZettelCheckedOut

	options := _ZettelsCheckinOptions{
		IncludeAkte: c.IncludeAkte,
		Format:      _ZettelFormatsText{},
	}

	if czs, err = zs.Checkout(options, args...); err != nil {
		err = _Error(err)
		return
	}

	files := make([]string, 0, len(czs))
	akten := make([]string, 0)

	for _, z := range czs {
		files = append(files, z.External.Path)

		if z.External.AktePath != "" {
			akten = append(akten, z.External.AktePath)
		}
	}

	if len(akten) > 0 {
		if err = _OpenFiles(akten...); err != nil {
			err = _Errorf("%q: %w", akten, err)
			return
		}
	}

	vimArgs := []string{
		"-c",
		"set ft=zit.zettel",
		"-c",
		"source ~/.vim/syntax/zit.zettel.vim",
	}

	if err = _OpenVimWithArgs(vimArgs, files...); err != nil {
		err = _Error(err)
		return
	}

	if _, err = zs.Checkin(options, files...); err != nil {
		err = _Error(err)
		return
	}

	return
}
