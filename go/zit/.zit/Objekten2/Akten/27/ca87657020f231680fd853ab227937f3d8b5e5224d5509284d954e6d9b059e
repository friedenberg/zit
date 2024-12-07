package files

import (
	"io"
	"os"
	"os/exec"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

const (
	buffer = 10
	limit  = 256
)

// TODO-P5 add exponential backoff for too many files open error
// TODO-P3 move away from openFilesGuard and honor too many files open error with
// exponential backoffs instead
type openFilesGuard struct {
	channel chan struct{}
}

var openFilesGuardInstance *openFilesGuard

func init() {
	// limitCmd := exec.Command("ulimit", "-S", "-n")
	// output, err := limitCmd.Output()

	// if err != nil {
	// 	panic(err)
	// }

	// limitStr := strings.TrimSpace(string(output))

	// limit, err := strconv.ParseInt(string(limitStr), 10, 64)

	// if err != nil {
	// 	panic(err)
	// }

	openFilesGuardInstance = &openFilesGuard{
		channel: make(chan struct{}, limit-buffer),
	}

	close(openFilesGuardInstance.channel)
}

func Len() int {
	return len(openFilesGuardInstance.channel)
}

func (g *openFilesGuard) Lock() {
	// g.channel <- struct{}{}
}

func (g *openFilesGuard) LockN(n int) {
	for i := 0; i < n; i++ {
		g.Lock()
	}
}

func (g *openFilesGuard) Unlock() {
	<-g.channel
}

func (g *openFilesGuard) UnlockN(n int) {
	for i := 0; i < n; i++ {
		g.Unlock()
	}
}

func CreateExclusiveWriteOnly(p string) (f *os.File, err error) {
	if f, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func CreateExclusiveWriteOnlyAndMaybeMakeDir(p string) (f *os.File, err error) {
	if f, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666); err != nil {
		if errors.IsNotExist(err) {
			dir := path.Dir(p)

			if err = os.MkdirAll(dir, os.ModeDir|0o755); err != nil {
				err = errors.Wrap(err)
				return
			}

			return CreateExclusiveWriteOnly(p)
		}

		err = errors.Wrap(err)
		return
	}

	return
}

func Create(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.Create(s); err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}

func OpenFile(name string, flag int, perm os.FileMode) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.OpenFile(name, flag, perm); err != nil {
		err = errors.Wrapf(err, "Mode: %d, Perm: %d", flag, perm)
		openFilesGuardInstance.Unlock()
	}

	return
}

func Open(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.Open(s); err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}

func OpenReadWrite(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.OpenFile(
		s,
		os.O_RDWR|os.O_CREATE,
		0o666,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func OpenCreateWriteOnlyTruncate(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.OpenFile(
		s,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0o666,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func OpenExclusiveWriteOnlyTruncate(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.OpenFile(
		s,
		os.O_WRONLY|os.O_EXCL|os.O_TRUNC,
		0o666,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func OpenExclusiveReadOnly(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.OpenFile(s, os.O_RDONLY|os.O_EXCL, 0o666); err != nil {
		err = errors.Wrapf(err, "Path: %q", s)
		return
	}

	return
}

func OpenExclusiveWriteOnly(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = os.OpenFile(s, os.O_WRONLY|os.O_EXCL, 0o666); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func Close(f *os.File) error {
	defer openFilesGuardInstance.Unlock()
	return f.Close()
}

func CombinedOutput(c *exec.Cmd) ([]byte, error) {
	openFilesGuardInstance.LockN(3)
	defer openFilesGuardInstance.UnlockN(3)

	return c.CombinedOutput()
}

func ReadAllString(s ...string) (o string, err error) {
	var f *os.File

	if f, err = Open(path.Join(s...)); err != nil {
		return
	}

	defer Close(f)

	var b []byte

	if b, err = io.ReadAll(f); err != nil {
		return
	}

	o = string(b)

	return
}
