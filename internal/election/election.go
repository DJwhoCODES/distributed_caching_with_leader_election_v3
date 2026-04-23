package election

import (
	"math/rand"
	"sync"
	"time"

	"github.com/djwhocodes/d-cache/internal/cluster"
)

type Election struct {
	mu       sync.Mutex
	state    State
	nodeID   string
	term     int
	votedFor string

	cluster *cluster.Cluster

	resetCh chan struct{}
}

func NewElection(nodeID string, c *cluster.Cluster) *Election {
	return &Election{
		state:   Follower,
		nodeID:  nodeID,
		cluster: c,
		resetCh: make(chan struct{}, 1),
	}
}

func (e *Election) Start() {
	go e.loop()
}

func (e *Election) loop() {
	for {
		timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond

		select {
		case <-time.After(timeout):
			e.startElection()
		case <-e.resetCh:
			// heartbeat received → reset timer
		}
	}
}

func (e *Election) startElection() {
	e.mu.Lock()
	e.state = Candidate
	e.term++
	e.votedFor = e.nodeID
	e.mu.Unlock()

	votes := 1
	nodes := e.cluster.GetNodes()

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, n := range nodes {
		if n.ID == e.nodeID {
			continue
		}

		wg.Add(1)

		go func(n *cluster.Node) {
			defer wg.Done()

			granted := e.requestVote(n.Addr)

			if granted {
				mu.Lock()
				votes++
				mu.Unlock()
			}
		}(n)
	}

	wg.Wait()

	if votes > len(nodes)/2 {
		e.becomeLeader()
	}
}

func (e *Election) becomeLeader() {
	e.mu.Lock()
	e.state = Leader
	e.mu.Unlock()

	e.cluster.SetLeader(e.nodeID)

	go e.startHeartbeat()
}

func (e *Election) startHeartbeat() {
	ticker := time.NewTicker(50 * time.Millisecond)

	for range ticker.C {
		if e.getState() != Leader {
			return
		}

		nodes := e.cluster.GetNodes()

		for _, n := range nodes {
			if n.ID == e.nodeID {
				continue
			}

			go e.sendHeartbeat(n.ID)
		}
	}
}

func (e *Election) getState() State {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.state
}

func (e *Election) OnVoteRequest(term int, candidateID string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if term < e.term {
		return false
	}

	e.term = term
	e.votedFor = candidateID
	e.state = Follower

	e.reset()

	return true
}

func (e *Election) OnHeartbeat(term int, leaderID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if term >= e.term {
		e.term = term
		e.state = Follower
		e.cluster.SetLeader(leaderID)
		e.reset()
	}
}

func (e *Election) reset() {
	select {
	case e.resetCh <- struct{}{}:
	default:
	}
}
