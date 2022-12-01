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



	chooseRunner()

	// var newmsg = "" 
	// for _,val := range runnersColor {
	// 	newmsg += string(val[1:]) + " "
	// }
	
	// err = writeToClients(clientsPresents, newmsg)

	// time.Sleep(1 * time.Second)

	for {
		checkArrival()

		handleResults()
	}

}


func chooseRunner() {
		
	runnerschose := [2]bool{false,false} //A changer pour 4 plus tard
	runnersColor := [2]string{}
	
	for {

		for i,client := range clientsPresents {
			if string(<-client.receiveChannel)[:2] == "3"+strconv.Itoa(i) {
				runnerschose[i] = true
				runnersColor[i] = string(<-client.receiveChannel)
			}
		}

		log.Println(runnerschose)
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
} 


func checkArrival() {
	ClientsFinished := [2]bool{false,false}
	for {
		for i,client := range clientsPresents {	
			if string(<-client.receiveChannel) == "50"+strconv.Itoa(i) {	
				ClientsFinished[i] = true
			}
		}
		c := 0
		for i := range ClientsFinished {
			if ClientsFinished[i] == true {
				c++
			}
		}
		if c == len(ClientsFinished) {
			writeToClients(clientsPresents,"600")
			break
		}
	}
}

func handleResults() {
	ClientsWantToRestart := [2]bool{false,false}
	for {
		for i,client := range clientsPresents {	
			if string(<-client.receiveChannel) == "70"+strconv.Itoa(i) {	
				ClientsWantToRestart[i] = true
			}
		}
		c := 0
		for i := range ClientsWantToRestart {
			if ClientsWantToRestart[i] == true {
				c++
			}
		}
		if c == len(ClientsWantToRestart) {
			writeToClients(clientsPresents,"800")
			break
		}
	}
}
























func writeMessage(client ClientListener, message string) (data int, err error) {
	data, err = client.conn.Write([]byte(message+"\n"))
	log.Println("Le message envoyé est : "+message)
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