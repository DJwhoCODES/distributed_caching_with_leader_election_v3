package protocol

import (
	"encoding/binary"
	"errors"
)

var (
	ErrIncompleteData = errors.New("incomplete data")
	ErrInvalidPacket  = errors.New("invalid packet")
)

func DecodeHeader(data []byte) (Header, error) {
	if len(data) < HeaderSize {
		return Header{}, ErrIncompleteData
	}

	h := Header{}
	h.Version = data[0]
	h.Command = Command(data[1])
	h.Flags = data[2]
	h.RequestID = binary.BigEndian.Uint32(data[3:7])
	h.PayloadLength = binary.BigEndian.Uint32(data[7:11])

	return h, nil
}

func DecodeRequest(data []byte) (*Request, error) {
	if len(data) < HeaderSize {
		return nil, ErrIncompleteData
	}

	header, err := DecodeHeader(data[:HeaderSize])

	if err != nil {
		return nil, err
	}

	totalLen := HeaderSize + int(header.PayloadLength)

	if len(data) < totalLen {
		return nil, ErrIncompleteData
	}

	body := data[HeaderSize:totalLen]

	offset := 0

	if offset+4 > len(body) {
		return nil, ErrInvalidPacket
	}

	keyLen := int(binary.BigEndian.Uint32(body[offset : offset+4]))
	offset += 4

	if offset+keyLen > len(body) {
		return nil, ErrInvalidPacket
	}

	key := body[offset : offset+keyLen]
	offset += keyLen

	if offset+4 > len(body) {
		return nil, ErrInvalidPacket
	}
	valLen := int(binary.BigEndian.Uint32(body[offset : offset+4]))
	offset += 4

	if offset+valLen > len(body) {
		return nil, ErrInvalidPacket
	}
	value := body[offset : offset+valLen]
	offset += valLen

	if offset+4 > len(body) {
		return nil, ErrInvalidPacket
	}
	ttl := binary.BigEndian.Uint32(body[offset : offset+4])

	req := &Request{
		Header: header,
		Key:    key,
		Value:  value,
		TTL:    ttl,
	}

	return req, nil
}

type StreamDecoder struct {
	buffer []byte
}

func NewStreamDecoder() *StreamDecoder {
	return &StreamDecoder{
		buffer: make([]byte, 0),
	}
}

func (d *StreamDecoder) Feed(data []byte) {
	d.buffer = append(d.buffer, data...)
}

func (d *StreamDecoder) Next() (*Request, error) {
	if len(d.buffer) < HeaderSize {
		return nil, ErrIncompleteData
	}

	header, err := DecodeHeader(d.buffer[:HeaderSize])
	if err != nil {
		return nil, err
	}

	totalLen := HeaderSize + int(header.PayloadLength)
	if len(d.buffer) < totalLen {
		return nil, ErrIncompleteData
	}

	packet := d.buffer[:totalLen]

	req, err := DecodeRequest(packet)
	if err != nil {
		return nil, err
	}

	d.buffer = d.buffer[totalLen:]

	return req, nil
}
