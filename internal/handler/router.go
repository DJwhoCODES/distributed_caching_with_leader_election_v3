package handler

import (
	"github.com/djwhocodes/d-cache/internal/cache"
	"github.com/djwhocodes/d-cache/internal/cluster"
	"github.com/djwhocodes/d-cache/internal/election"
	"github.com/djwhocodes/d-cache/internal/protocol"
	"github.com/djwhocodes/d-cache/internal/replication"
)

type Router struct {
	store      *cache.Store
	cluster    *cluster.Cluster
	replicator *replication.Replicator
	election   *election.Election
}

func NewRouter(store *cache.Store, c *cluster.Cluster) *Router {
	return &Router{
		store:      store,
		cluster:    c,
		replicator: replication.NewReplicator(c),
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

	case protocol.CmdVoteRequest:
		return r.handleVote(req)

	case protocol.CmdHeartbeat:
		return r.handleHeartbeat(req)

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

func (r *Router) handleVote(req *protocol.Request) *protocol.Response {
	candidateID := string(req.Value)
	term := int(req.TTL)

	granted := r.election.OnVoteRequest(term, candidateID)

	status := protocol.StatusError
	if granted {
		status = protocol.StatusOK
	}

	return &protocol.Response{
		Header: protocol.Header{
			Command:   req.Header.Command,
			RequestID: req.Header.RequestID,
		},
		Status: status,
	}
}

func (r *Router) handleHeartbeat(req *protocol.Request) *protocol.Response {
	leaderID := string(req.Value)
	term := int(req.TTL)

	r.election.OnHeartbeat(term, leaderID)

	return &protocol.Response{
		Header: protocol.Header{
			Command:   req.Header.Command,
			RequestID: req.Header.RequestID,
		},
		Status: protocol.StatusOK,
	}
}
