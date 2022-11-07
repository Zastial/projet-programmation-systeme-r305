package main

import (
	"log"
	"net"
	"bufio"
)

func connexion(ip string, connexion []net.Conn) (clients []net.Conn) {
	conn, err := net.Dial("tcp", ip + ":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	clients = append(connexion, conn)

	return clients
}

func writeToServer(clients []net.Conn, msg string) {
	for _,client := range clients {
		_, err = client.Write([]byte(msg))
	}
}

func readFromServer(clients) (text string){
	for {
		for _,conn := range clients {
			netData, _ := bufio.NewReader(conn).ReadString('\n')
			text = string(netData)
		}
		return text
	}
}