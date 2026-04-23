package handler

import "github.com/djwhocodes/d-cache/internal/protocol"

func (r *Router) handleGet(req *protocol.Request) *protocol.Response {
	val, ok := r.store.Get(string(req.Key))

	if !ok {
		return &protocol.Response{
			Header: protocol.Header{
				Command:   req.Header.Command,
				RequestID: req.Header.RequestID,
			},
			Status: protocol.StatusNotFound,
		}
	}

	return &protocol.Response{
		Header: protocol.Header{
			Command:   req.Header.Command,
			RequestID: req.Header.RequestID,
		},
		Status: protocol.StatusOK,
		Value:  val,
	}
}
