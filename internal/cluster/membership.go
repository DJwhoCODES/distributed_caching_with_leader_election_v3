package cluster

import "sync"

type Cluster struct {
	mu       sync.RWMutex
	nodes    map[string]*Node
	self     *Node
	leaderID string
}

func NewCluster(self *Node) *Cluster {
	return &Cluster{
		nodes: make(map[string]*Node),
		self:  self,
	}
}

func (c *Cluster) AddNode(n *Node) {
	c.mu.Lock()
	c.nodes[n.ID] = n
	c.mu.Unlock()
}

func (c *Cluster) GetNodes() []*Node {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]*Node, 0, len(c.nodes))
	for _, n := range c.nodes {
		out = append(out, n)
	}
	return out
}

func (c *Cluster) SetLeader(id string) {
	c.mu.Lock()
	c.leaderID = id
	c.mu.Unlock()
}

func (c *Cluster) GetLeader() *Node {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nodes[c.leaderID]
}
