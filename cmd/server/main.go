package main

import (
	"log"
	"net"

	"github.com/djwhocodes/d-cache/internal/protocol"
	"github.com/djwhocodes/d-cache/internal/transport"
)

func main() {
	server := transport.NewTCPServer(":8080")

	err := server.Start(connectionHandler)

	if err != nil {
		log.Fatal(err)
	}
}

func connectionHandler(conn net.Conn) {
	transport.HandleConnection(conn, handleRequest)
}

func handleRequest(req *protocol.Request) *protocol.Response {
	return &protocol.Response{
		Header: protocol.Header{
			Command:   req.Header.Command,
			RequestID: req.Header.RequestID,
		},
		Status: protocol.StatusOK,
		Value:  []byte("OK"),
	}
}
