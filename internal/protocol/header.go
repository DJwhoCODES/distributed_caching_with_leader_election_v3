package protocol

type Header struct {
	Version       uint8
	Command       Command
	Flags         uint8
	RequestID     uint32
	PayloadLength uint32
}

const HeaderSize = 11
