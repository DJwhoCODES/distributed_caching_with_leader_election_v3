package protocol

import "encoding/binary"

func DecodeResponse(data []byte) (*Response, error) {
	if len(data) < HeaderSize {
		return nil, ErrIncompleteData
	}

	header, err := DecodeHeader(data[:HeaderSize])
	if err != nil {
		return nil, err
	}

	body := data[HeaderSize:]
	offset := 0

	// Status
	status := Status(body[offset])
	offset += 1

	// Value
	valLen := int(binary.BigEndian.Uint32(body[offset : offset+4]))
	offset += 4

	value := body[offset : offset+valLen]
	offset += valLen

	// Error
	errLen := int(binary.BigEndian.Uint32(body[offset : offset+4]))
	offset += 4

	errMsg := string(body[offset : offset+errLen])

	return &Response{
		Header: header,
		Status: status,
		Value:  value,
		Error:  errMsg,
	}, nil
}
