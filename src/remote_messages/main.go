package remote_messages

type MessageType int

const (
	MessageTypeUnknown = MessageType(iota)
	MessageTypeSenderHi
	MessageTypeReceiverHi
	MessageTypeSenderPushTransaction
	MessageTypeReceiverWantObjektenAndAkten
	MessageTypeSenderShareObjektenAndAktenSockets
	MessageTypeReceiverDone
)

type MessageContents interface{}

type Message struct {
	MessageType
	Contents MessageContents
}

func (m *Message) NextLine() (ok bool) {
	if m.MessageType == MessageTypeReceiverDone {
		return
	}

	ok = true
	m.MessageType += 1

	return
}
