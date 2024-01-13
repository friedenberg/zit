package remote_conn

type Script interface {
	HandleSenderDialogue(Dialogue) error
	HandleReceiverDialogue(Dialogue) error
}
