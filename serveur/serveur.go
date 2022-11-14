package main

import (
	"fmt"
	"log"
	"net"

)

var connexion []net.Conn

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		connexion = append(connexion, conn)
		defer conn.Close()

		fmt.Println("Un client a rejoint le serveur")



		var msg = make([]byte, 1024)
		byteCount, err := conn.Read(msg)
		if err != nil {
			log.Println("error", err)
			return
		}
		log.Println("Bits reçu:", byteCount)
		log.Println("Message reçu:", string(msg))
	}
}

func giveConn() (connexion[]net.Conn) {
	return connexion
}

