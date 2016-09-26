package wire

type MsgBroadcast struct {
	content string
}

func (m *MsgBroadcast) Command() string {
	return CmdBroadcast
}

func NewBroadcastMsg(content string) Message {
	return &MsgBroadcast{content}
}
