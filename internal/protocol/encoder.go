package protocol

import (
	"bytes"
	"encoding/binary"
)

// EncodeRequest -> []byte (wire format)
func EncodeRequest(req *Request) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	body := bytes.NewBuffer(make([]byte, 0))

	if err := binary.Write(body, binary.BigEndian, uint32(len(req.Key))); err != nil {
		return nil, err
	}

	body.Write(req.Key)

	if err := binary.Write(body, binary.BigEndian, uint32(len(req.Value))); err != nil {
		return nil, err
	}

	body.Write(req.Value)

	if err := binary.Write(body, binary.BigEndian, req.TTL); err != nil {
		return nil, err
	}

	req.Header.Version = Version
	req.Header.PayloadLength = uint32(body.Len())

	if err := binary.Write(buf, binary.BigEndian, req.Header.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.Header.Command); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.Header.Flags); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.Header.RequestID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.Header.PayloadLength); err != nil {
		return nil, err
	}

	buf.Write(body.Bytes())

	return buf.Bytes(), nil
}

func EncodeResponse(res *Response) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	body := bytes.NewBuffer(make([]byte, 0))

	if err := binary.Write(body, binary.BigEndian, res.Status); err != nil {
		return nil, err
	}

	if err := binary.Write(body, binary.BigEndian, uint32(len(res.Value))); err != nil {
		return nil, err
	}

	body.Write(res.Value)

	errBytes := []byte(res.Error)
	if err := binary.Write(body, binary.BigEndian, uint32(len(errBytes))); err != nil {
		return nil, err
	}

	body.Write(errBytes)

	res.Header.Version = Version
	res.Header.PayloadLength = uint32(body.Len())

	binary.Write(buf, binary.BigEndian, res.Header.Version)
	binary.Write(buf, binary.BigEndian, res.Header.Command)
	binary.Write(buf, binary.BigEndian, res.Header.Flags)
	binary.Write(buf, binary.BigEndian, res.Header.RequestID)
	binary.Write(buf, binary.BigEndian, res.Header.PayloadLength)

	buf.Write(body.Bytes())

	return buf.Bytes(), nil
}
