package replication

import (
	"sync"

	"github.com/djwhocodes/d-cache/internal/client"
	"github.com/djwhocodes/d-cache/internal/cluster"
	"github.com/djwhocodes/d-cache/internal/protocol"
)

type Replicator struct {
	cluster *cluster.Cluster
}

func NewReplicator(c *cluster.Cluster) *Replicator {
	return &Replicator{
		cluster: c,
	}
}

func (r *Replicator) Replicate(req *protocol.Request) {
	nodes := r.cluster.GetNodes()

	var wg sync.WaitGroup

	for _, n := range nodes {
		// skip self (leader)
		if n.ID == r.cluster.GetLeader().ID {
			continue
		}

		wg.Add(1)

		go func(addr string) {
			defer wg.Done()

			c := client.NewTCPClient(addr)

			// mark as internal replication
			req.Header.Flags = 1

			_, _ = c.Send(req) // ignore error for now
		}(n.Addr)
	}

	wg.Wait()
}
