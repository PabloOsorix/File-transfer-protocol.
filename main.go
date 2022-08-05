package main

import (
	//"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8031")
	if err != nil {
		log.Printf("%v", err)
	}

	hub := newHub()
	chargeFiles(*hub)
	go hub.run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("%v", err)
		}

		c := newClient(
			conn,
			hub.commands,
			hub.registrations,
			hub.deregistrations,
		)
		go c.read()
	}
}
