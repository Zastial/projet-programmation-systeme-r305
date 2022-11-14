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

		var c = 0
		handleWelcome := make([]bool, len(connexion))
		for (c != 4) {
			for i, conn := range connexion {
				if (HandleWelcomeScreen() == true) {
					log.Println("zizi")
					// handleWelcome[i] = true
					// c++
				}
			}
		}

		if(len(connexion) == 4) {
			
		}
	}
}

