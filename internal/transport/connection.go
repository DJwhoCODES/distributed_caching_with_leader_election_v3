package transport

import (
	"io"
	"log"
	"net"

	"github.com/djwhocodes/d-cache/internal/protocol"
)

const ReadBufferSize = 4096

func HandleConnection(conn net.Conn, requestHandler func(*protocol.Request) *protocol.Response) {
	defer conn.Close()

	decoder := protocol.NewStreamDecoder()
	buf := make([]byte, ReadBufferSize)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			return
		}

		decoder.Feed(buf[:n])

		for {
			req, err := decoder.Next()
			if err != nil {
				if err == protocol.ErrIncompleteData {
					break
				}
				log.Println("decode error:", err)
				return
			}

			res := requestHandler(req)

			bytes, err := protocol.EncodeResponse(res)

			if err != nil {
				log.Println("encode error:", err)
				return
			}

			_, err = conn.Write(bytes)
			if err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}
}
