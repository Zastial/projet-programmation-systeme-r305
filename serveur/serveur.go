package main

import (
	"fmt"
	"log"
	"net"
	"bufio"
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
		for(len(connexion) < 4) {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("accept error:", err)
				return
			}
			connexion = append(connexion, conn)
			defer conn.Close()

			fmt.Println("Un client a rejoint le serveur")
		}
	}

	if (len(connexion) == 4) {
		for _, conn := range connexion {
			writer := bufio.NewWriter(conn)
			writer.WriteString("200\n")
			writer.Flush()
			if err != nil{
				return
			}
		}
	} 

	for {

	}
}
