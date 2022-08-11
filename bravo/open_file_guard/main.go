package open_file_guard

import (
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const (
	buffer = 35
)

//TODO-P5 add exponential backoff for too many files open error
type openFilesGuard struct {
	channel chan struct{}
}

var openFilesGuardInstance *openFilesGuard

func init() {
	limitCmd := exec.Command("ulimit", "-S", "-n")
	output, err := limitCmd.Output()

	if err != nil {
		panic(err)
	}

	limitStr := strings.TrimSpace(string(output))

	limit, err := strconv.ParseInt(string(limitStr), 10, 64)

	if err != nil {
		panic(err)
	}

	openFilesGuardInstance = &openFilesGuard{
		channel: make(chan struct{}, limit-buffer),
	}
}

func Len() int {
	return len(openFilesGuardInstance.channel)
}

func (g *openFilesGuard) Lock() {
	g.channel <- struct{}{}
	// logz.Caller(3, "locked: %d", len(g.channel))
}

func (g *openFilesGuard) LockN(n int) {
	for i := 0; i < n; i++ {
		g.Lock()
	}
}

func (g *openFilesGuard) Unlock() {
	<-g.channel
	// logz.Caller(3, "unlocked %d", len(g.channel))
}

func (g *openFilesGuard) UnlockN(n int) {
	for i := 0; i < n; i++ {
		g.Unlock()
	}
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
