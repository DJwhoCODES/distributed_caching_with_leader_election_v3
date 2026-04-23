package handler

import (
	"github.com/djwhocodes/d-cache/internal/cache"
	"github.com/djwhocodes/d-cache/internal/protocol"
)

type Router struct {
	store *cache.Store
}

func NewRouter(store *cache.Store) *Router {
	return &Router{
		store: store,
	}
}

func (r *Router) Handle(req *protocol.Request) *protocol.Response {
	switch req.Header.Command {
	case protocol.CmdGet:
		return r.handleGet(req)

	case protocol.CmdSet:
		return r.handleSet(req)

	case protocol.CmdDelete:
		return r.handleDelete(req)

	case protocol.CmdPing:
		return r.handlePing(req)

	default:
		return &protocol.Response{
			Header: protocol.Header{
				Command:   req.Header.Command,
				RequestID: req.Header.RequestID,
			},
			Status: protocol.StatusError,
			Error:  "unknown command",
		}
	}
}
