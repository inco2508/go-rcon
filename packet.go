package rcon

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/rand"
	"net"
)

type Packet struct {
	Size        int32
	ID          int32
	Type        int32
	Body        string
	EmptyString []byte
}

func newPacket(packetType int32, body string) *Packet {
	return &Packet{
		ID:          int32(rand.Int31()),
		Type:        packetType,
		Size:        int32(4 + 4 + len(body) + 2), // ID + Type + Body + EmptyString
		Body:        body,
		EmptyString: []byte{0, 0},
	}
}

func (p Packet) encode() ([]byte, error) {
	var buffer = bytes.Buffer{}

	binary.Write(&buffer, binary.LittleEndian, uint32(p.Size))
	binary.Write(&buffer, binary.LittleEndian, uint32(p.ID))
	binary.Write(&buffer, binary.LittleEndian, uint32(p.Type))
	buffer.Write([]byte(p.Body))
	buffer.Write(p.EmptyString)

	return buffer.Bytes(), nil
}

func (p *Packet) decode(data []byte) error {
	p.Size = int32(binary.LittleEndian.Uint32(data[0:4]))
	p.ID = int32(binary.LittleEndian.Uint32(data[4:8]))
	p.Type = int32(binary.LittleEndian.Uint32(data[8:12]))
	p.Body = string(data[12 : len(data)-2])
	p.EmptyString = data[len(data)-2:]

	return nil
}

func (p *Packet) Write(conn net.Conn) (int, error) {
	encoded, _ := p.encode()
	n, err := conn.Write(encoded)
	if err != nil {
		return 0, err
	}

	return n, nil

}

func (p *Packet) Read(conn net.Conn) (int, error) {
	decoded := &Packet{}
	readBuffer := make([]byte, 4096)
	n, err := conn.Read(readBuffer)
	if err != nil {
		if err != io.EOF {
			return 0, err
		}
	}

	decoded.decode(readBuffer)

	return n, nil
}

func NewAuthPacket(password string) *Packet {
	return newPacket(3, password)
}

func NewCommandPacket(command string) *Packet {
	return newPacket(2, command)
}
