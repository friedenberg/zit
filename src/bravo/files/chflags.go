package files

import (
	"io/fs"
	"os/exec"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func isNotExists(err error, msg []byte) bool {
	return strings.HasSuffix(string(msg), "No such file or directory")
}

func setUserChanges(paths []string, allow bool) (err error) {
	setting := "uchg"

	if allow {
		setting = "no" + setting
	}

	// TODO-P2 change to syscall:
	// https://github.com/snapcore/snapd/blob/master/osutil/chattr.go
	// https://stackoverflow.com/questions/69542185/make-file-immutable-syscall-chflagsfilename
	cmd := exec.Command(
		"/usr/bin/chflags",
		append([]string{setting}, paths...)...,
	)

	var msg []byte

	msg, err = cmd.CombinedOutput()

	if err != nil {
		if isNotExists(err, msg) {
			err = fs.ErrNotExist
		} else {
			err = errors.Errorf("failed to run chflags: %q", msg)
		}

		return
	}

	return
}

func SetAllowUserChanges(paths ...string) (err error) {
	return setUserChanges(paths, true)
}

func SetDisallowUserChanges(paths ...string) (err error) {
	return setUserChanges(paths, false)
}
