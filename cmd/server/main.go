package main

import (
	"log"
	"net"
	"time"

	"github.com/djwhocodes/d-cache/internal/cache"
	"github.com/djwhocodes/d-cache/internal/cluster"
	"github.com/djwhocodes/d-cache/internal/election"
	"github.com/djwhocodes/d-cache/internal/handler"
	"github.com/djwhocodes/d-cache/internal/transport"
)

func main() {
	server := transport.NewTCPServer(":8080")
	store := cache.NewStore()

	node := &cluster.Node{
		ID:   "node-1",
		Addr: ":8080",
	}

	cluster := cluster.NewCluster(node)
	cluster.AddNode(node)

	store.StartJanitor(10 * time.Second)

	e := election.NewElection(node.ID, cluster)
	e.Start()

	router := handler.NewRouter(store, cluster)
	router.election = e

	err := server.Start(connectionHandler, *router)

	if err != nil {
		log.Fatal(err)
	}
}

func connectionHandler(conn net.Conn, router handler.Router) {
	transport.HandleConnection(conn, router.Handle)
}
