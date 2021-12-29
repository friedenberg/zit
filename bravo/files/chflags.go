package files

import (
	"io/fs"
	"os/exec"
	"strings"
)

func isNotExists(err error, msg []byte) bool {
	return strings.HasSuffix(string(msg), "No such file or directory")
}

func setUserChanges(paths []string, allow bool) (err error) {
	setting := "uchg"

	if allow {
		setting = "no" + setting
	}

	cmd := exec.Command("chflags", append([]string{setting}, paths...)...)
	var msg []byte
	msg, err = cmd.CombinedOutput()

	if err != nil {
		if isNotExists(err, msg) {
			err = fs.ErrNotExist
		} else {
			err = _Errorf("failed to run chflags: %s", msg)
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
