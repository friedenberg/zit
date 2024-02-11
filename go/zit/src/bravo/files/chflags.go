package files

import (
	"io/fs"
	"os/exec"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

func isNotExists(err error, msg []byte) bool {
	return strings.HasSuffix(string(msg), "No such file or directory")
}

type userChangesOptions struct {
	allow     bool
	recursive bool
}

func setUserChanges(paths []string, options userChangesOptions) (err error) {
	var args []string

	if options.recursive {
		args = append(args, "-R")
	}

	setting := "uchg"

	if options.allow {
		setting = "no" + setting
	}

	args = append(args, setting)

	// TODO-P2 change to syscall:
	// https://github.com/snapcore/snapd/blob/master/osutil/chattr.go
	// https://stackoverflow.com/questions/69542185/make-file-immutable-syscall-chflagsfilename
	cmd := exec.Command(
		"/usr/bin/chflags",
		append(args, paths...)...,
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

func SetAllowUserChangesRecursive(paths ...string) (err error) {
	return setUserChanges(
		paths,
		userChangesOptions{allow: true, recursive: true},
	)
}

func SetAllowUserChanges(paths ...string) (err error) {
	return setUserChanges(paths, userChangesOptions{allow: true})
}

func SetDisallowUserChanges(paths ...string) (err error) {
	return setUserChanges(paths, userChangesOptions{allow: false})
}
