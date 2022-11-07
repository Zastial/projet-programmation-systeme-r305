package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	var connexion []net.Conn
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		connexion = append(connexion, conn)
		defer conn.Close()

		fmt.Println("Un client a rejoint le serveur")

		if(len(connexion) == 4) {
			return
		}
	}
}
