package main

import (
	"log"
	"net"
	"time"

	"github.com/djwhocodes/d-cache/internal/cache"
	"github.com/djwhocodes/d-cache/internal/handler"
	"github.com/djwhocodes/d-cache/internal/transport"
)

func main() {
	server := transport.NewTCPServer(":8080")
	store := cache.NewStore()
	store.StartJanitor(10 * time.Second)

	router := handler.NewRouter(store)

	err := server.Start(connectionHandler, *router)

	if err != nil {
		log.Fatal(err)
	}
}

func connectionHandler(conn net.Conn, router handler.Router) {
	transport.HandleConnection(conn, router.Handle)
}
