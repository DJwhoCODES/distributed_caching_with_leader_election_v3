package client

import (
	"net"
	"time"

	"github.com/djwhocodes/d-cache/internal/protocol"
)

type TCPClient struct {
	Addr string
}

func NewTCPClient(addr string) *TCPClient {
	return &TCPClient{Addr: addr}
}

func (c *TCPClient) Send(req *protocol.Request) (*protocol.Response, error) {
	conn, err := net.DialTimeout("tcp", c.Addr, 2*time.Second)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	data, err := protocol.EncodeRequest(req)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(data)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)

	if err != nil {
		return nil, err
	}

	return protocol.DecodeResponse(buf[:n])
}
