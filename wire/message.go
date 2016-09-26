package wire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io"
)

type NetID uint32

const (
	HeaderLen = 24

	CommandSize = 12

	MainNetID NetID = 0xdadb1986
	TestNetID NetID = 0xccdd2086

	CmdVersion   = "version"
	CmdBroadcast = "broadcast"
)

var (
	byteOrder = binary.LittleEndian
)

type Message interface {
	Command() string
}

type MessageHeader struct {
	magic    NetID
	command  string
	msglen   uint32
	checksum [4]byte
}

func ReadMessage(r io.Reader, magic NetID) (Message, error) {
	var buf [HeaderLen]byte

	n, err := io.ReadFull(r, buf[:])
	if n != HeaderLen || err != nil {
		return nil, err
	}

	hr := bytes.NewReader(buf[:])

	var header MessageHeader
	var command [CommandSize]byte

	readElements(hr, &header.magic, &command, &header.msglen, &header.checksum)

	header.command = string(bytes.TrimRight(command[:], string(0)))

	if header.magic != magic {
		return nil, fmt.Errorf("Message Header Magic: %v doesn't match server magic: %v", header.magic, magic)
	}

	msg, err := getDefaultMsg(header.command)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, header.msglen)
	_, err = io.ReadFull(r, payload[:])
	if err != nil {
		return nil, err
	}

	err = msgpack.Unmarshal(payload, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func WriteMessage(w io.Writer, magic NetID, msg Message) error {
	payload, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}

	size := len(payload)
	cmd := msg.Command()

	var command [CommandSize]byte
	copy(command[:], []byte(cmd))

	buf := new(bytes.Buffer)

	writeElements(buf, magic, command[:], size)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(payload)
	return err
}

func writeElements(w io.Writer, args ...interface{}) error {
	for _, el := range args {
		err := writeElement(w, el)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeElement(w io.Writer, el interface{}) error {
	switch e := el.(type) {
	case NetID:
		return writeUint32(w, uint32(e))
	case uint32:
		return writeUint32(w, e)
	case int:
		return writeUint32(w, uint32(e))
	case []byte:
		_, err := w.Write(e)
		return err
	}

	return binary.Write(w, byteOrder, el)
}

func writeUint32(w io.Writer, i uint32) error {
	var buf [4]byte

	byteOrder.PutUint32(buf[:], i)

	_, err := w.Write(buf[:])
	return err
}

func readElements(r io.Reader, args ...interface{}) error {
	for _, el := range args {
		err := readElement(r, el)
		if err != nil {
			return err
		}
	}

	return nil
}

func readElement(r io.Reader, el interface{}) error {
	switch e := el.(type) {
	case *uint32:
		v, err := readUint32(r)
		if err != nil {
			return err
		}

		*e = v
	case *NetID:
		v, err := readUint32(r)
		if err != nil {
			return err
		}

		*e = NetID(v)
	case *[12]byte:
		_, err := io.ReadFull(r, e[:])
		if err != nil {
			return err
		}

		return nil
	case *[4]byte:
		_, err := io.ReadFull(r, e[:])
		if err != nil {
			return err
		}

		return nil
	}

	return binary.Read(r, byteOrder, el)
}

func readUint32(r io.Reader) (uint32, error) {
	var buf [4]byte

	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}

	return byteOrder.Uint32(buf[:]), nil
}

func getDefaultMsg(command string) (Message, error) {
	switch command {
	case CmdVersion:
		return &MsgVersion{}, nil
	case CmdBroadcast:
		return &MsgBroadcast{}, nil
	default:
		return nil, fmt.Errorf("Invalid command: %s", command)
	}
}
