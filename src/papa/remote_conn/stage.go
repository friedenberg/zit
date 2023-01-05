package remote_conn

type stage struct {
	sockPath     string
	mainDialogue Dialogue
}

func (s stage) MainDialogue() Dialogue {
	return s.mainDialogue
}
