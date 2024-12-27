package repo

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

//go:generate stringer -type=RemoteType
type RemoteType int

const (
	RemoteTypeUnspecified = RemoteType(iota)
	RemoteTypeNativeDotenvXDG
	RemoteTypeSocketUnix
	RemoteTypePortHTTP
	RemoteTypeStdioLocal
	RemoteTypeStdioSSH
)

func (t *RemoteType) Set(v string) (err error) {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "", "none", "unspecified":
		*t = RemoteTypeUnspecified

	case "native-dotenv-xdg":
		*t = RemoteTypeNativeDotenvXDG

	case "socket-unix":
		*t = RemoteTypeSocketUnix

	case "port-http":
		*t = RemoteTypePortHTTP

	case "stdio-local":
		*t = RemoteTypeStdioLocal

	case "stdio-ssh":
		*t = RemoteTypeStdioSSH

	default:
		err = errors.Errorf("unsupported remote type: %q", v)
		return
	}

	return
}
