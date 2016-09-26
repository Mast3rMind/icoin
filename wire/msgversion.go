package wire

type MsgVersion struct {
}

func (m *MsgVersion) Command() string {
	return CmdVersion
}
