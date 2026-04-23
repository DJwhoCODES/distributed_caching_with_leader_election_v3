package transport

import (
	"log"
	"net"

	"github.com/djwhocodes/d-cache/internal/handler"
)

type TCPServer struct {
	Addr string
}

func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		Addr: addr,
	}
}

func (s *TCPServer) Start(
	handler func(net.Conn, handler.Router),
	router handler.Router,
) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	log.Println("TCP server listening on", s.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go handler(conn, router)
	}
}
