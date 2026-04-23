package handler

import "github.com/djwhocodes/d-cache/internal/protocol"

func (r *Router) handlePing(req *protocol.Request) *protocol.Response {
	return &protocol.Response{
		Header: protocol.Header{
			Command:   req.Header.Command,
			RequestID: req.Header.RequestID,
		},
		Status: protocol.StatusOK,
		Value:  []byte("PONG"),
	}
}
