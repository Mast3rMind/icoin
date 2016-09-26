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

	CmdVersion = "version"
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
	default:
		return nil, fmt.Errorf("Invalid command: %s", command)
	}
}
