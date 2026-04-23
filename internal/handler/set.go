package handler

import "github.com/djwhocodes/d-cache/internal/protocol"

func (r *Router) handleSet(req *protocol.Request) *protocol.Response {
	r.store.Set(string(req.Key), req.Value, req.TTL)

	if r.cluster.IsLeader() && req.Header.Flags == 0 {
		r.replicator.Replicate(req)
	}

	return &protocol.Response{
		Header: protocol.Header{
			Command:   req.Header.Command,
			RequestID: req.Header.RequestID,
		},
		Status: protocol.StatusOK,
	}
}
