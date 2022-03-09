package user_ops

import "github.com/friedenberg/zit/bravo/open_file_guard"

type OpenVim struct {
	Options []string
}

type OpenVimResults struct {
}

func (c OpenVim) Run(args ...string) (results OpenVimResults, err error) {
	vimArgs := make([]string, 0, len(c.Options)*2)

	for _, o := range c.Options {
		vimArgs = append(vimArgs, "-c", o)
	}

	if err = open_file_guard.OpenVimWithArgs(vimArgs, args...); err != nil {
		err = _Error(err)
		return
	}

	return
}
