package main

import (
	"log"
	"net"
	"bufio"
	"strings"
	"time"
)

type ClientListener struct {
	conn net.Conn
	receiveChannel chan string
	color int
}

var clientsPresents []ClientListener

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()


	for(len(clientsPresents) < 4) {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		defer conn.Close()


		var cl ClientListener
		cl.conn = conn
		cl.receiveChannel = make(chan string, 5)
		clientsPresents = append(clientsPresents,cl)


		err = writeMessage(cl, "100")
		if err != nil {
			log.Println("write error:", err)
			return
		}
		go receiveFromClient(cl)
		log.Println("Client connected")
	}

	err = writeToClients(clientsPresents,"200")
	if err != nil {
		log.Println("write error:", err)
		return
	}

	time.Sleep(10*time.Second)
}

func writeMessage(client ClientListener, message string) (err error) {
	_, err = client.conn.Write([]byte(message+"\n"))
	return err
}

func writeToClients(clients []ClientListener, message string) (err error) {
	for _, client := range clients {
		err = writeMessage(client,message)
		if err != nil {
			log.Println("listen error:", err)
			return
		}
		log.Println("message envoyé : ",message)
	}
	return err
}

func receiveFromClient(client ClientListener){
	reader := bufio.NewReader(client.conn)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		strip := strings.TrimSuffix(s, "\n")
		log.Println("received message from client : ", strip)

		client.receiveChannel <- strip
	}
}
