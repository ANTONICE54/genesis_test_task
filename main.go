package main

import (
	"genesis_tt/api"
	"log"
)

// In main function the HTTP server is initialized
func main() {

	server, err := api.NewServer()
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}

	go server.ListenForMail()
	go server.ListenForShutdown()

	err = server.Start()
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

}
