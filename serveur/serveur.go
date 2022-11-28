package main

import (
	"log"
	"net"
	"bufio"
	"strings"
	"strconv"
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

	var maxClients = 2

	for(len(clientsPresents) < maxClients) {
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


		data,err := writeMessage(cl, "100")
		log.Println(data)
		if err != nil {
			log.Println("write error:", err)
			return
		}

		data,err = writeMessage(cl,strconv.Itoa(len(clientsPresents)))
		log.Println("id : ",data)
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

	runnerschose := [2]bool{false,false} //A changer pour 4 plus tard
	runnersColor := [2]string{}
	for {
		for i,client := range clientsPresents {	
			go receiveFromClient(client)
			if string(<-client.receiveChannel)[:2] == "3"+strconv.Itoa(i) {	
				log.Println("OK CA MARCHE BB")
				runnerschose[i] = true
				runnersColor[i] = string(<-client.receiveChannel)
			}
		}
		c := 0
		for i := range runnerschose {
			if runnerschose[i] == true {
				c++
			}
		}
		if c == len(runnerschose) {
			writeToClients(clientsPresents,"400")
			break
		}
	}

	log.Println(runnerschose)
	log.Println(runnersColor)

	for {

	}

}

func writeMessage(client ClientListener, message string) (data int, err error) {
	data, err = client.conn.Write([]byte(message+"\n"))
	return data,err
}

func writeToClients(clients []ClientListener, message string) (err error) {
	for _, client := range clients {
		_,err := writeMessage(client,message)
		if err != nil {
			log.Println("listen error:", err)
			return err
		}
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
