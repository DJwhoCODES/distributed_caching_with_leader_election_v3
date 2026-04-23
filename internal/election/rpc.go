package election

import (
	"github.com/djwhocodes/d-cache/internal/client"
	"github.com/djwhocodes/d-cache/internal/protocol"
)

func (e *Election) requestVote(nodeAddr string) bool {
	c := client.NewTCPClient(nodeAddr)

	req := &protocol.Request{
		Header: protocol.Header{
			Command: protocol.CmdVoteRequest,
		},
		Value: []byte(e.nodeID), // candidateID
		TTL:   uint32(e.term),   // reuse TTL as term (simple hack)
	}

	res, err := c.Send(req)
	if err != nil {
		return false
	}

	return res.Status == protocol.StatusOK
}

func (e *Election) sendHeartbeat(nodeAddr string) {
	c := client.NewTCPClient(nodeAddr)

	req := &protocol.Request{
		Header: protocol.Header{
			Command: protocol.CmdHeartbeat,
		},
		Value: []byte(e.nodeID), // leaderID
		TTL:   uint32(e.term),
	}

	_, _ = c.Send(req)
}
