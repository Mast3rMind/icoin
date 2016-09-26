package wire

type MsgBroadcast struct {
	Content string
}

func (m *MsgBroadcast) Command() string {
	return CmdBroadcast
}

func NewBroadcastMsg(content string) Message {
	return &MsgBroadcast{content}
}
