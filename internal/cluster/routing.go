package cluster

import (
	"hash/fnv"
	"sort"
)

func (c *Cluster) GetNodeForKey(key string) *Node {
	nodes := c.GetNodes()

	if len(nodes) == 0 {
		return nil
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})

	hash := fnv.New32a()

	hash.Write([]byte(key))

	idx := int(hash.Sum32()) % len(nodes)

	return nodes[idx]
}

func (c *Cluster) IsLeader() bool {
	return c.self.ID == c.leaderID
}
