package remote_conn

import "net"

type stage struct {
	sockPath     string
	address      *net.UnixAddr
	mainDialogue Dialogue
}

func (s stage) MainDialogue() Dialogue {
	return s.mainDialogue
}
