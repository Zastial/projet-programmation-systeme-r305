package main

import (
	//"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	//"strings"
	//"time"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	var connexion []net.Conn
	for{
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		connexion = append(connexion, conn)
		defer conn.Close()

		fmt.Println("Un client a rejoint le serveur")

		conn.Write([]byte("Bienvenue sur le serveur\n"))

		for i := 0; i <= len(connexion)-1 ;i++ {
			connexion[i].Write([]byte("Joueurs : "+ strconv.Itoa(len(connexion)) +"/4\n"))
		}


		if(len(connexion) == 4) {
			for i := 0; i < len(connexion);i++ {
				connexion[i].Write([]byte("Le serveur est plein\n"))
			}
			return
		}
	}
}