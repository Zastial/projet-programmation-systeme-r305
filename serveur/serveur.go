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
	good bool
}

var clientsPresents []ClientListener

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	var maxClients = 4

	for(len(clientsPresents) < maxClients) {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		defer conn.Close()


		var cl ClientListener
		cl.conn = conn
		cl.receiveChannel = make(chan string, 100)
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


	for {
		checkArrival()

		handleResults()
	}

}

func playerSelector() {

	ancienMess0 := ""
	ancienMess1 := ""
	ancienMess2 := ""
	ancienMess3 := ""

	for {
		select {
		case mess := <-clientsPresents[0].receiveChannel:
			for i:=1;i<=8;i++ {
				if string(mess) == "4"+strconv.Itoa(0)+strconv.Itoa(i) && mess != ancienMess0 {
					writeToClients(clientsPresents,mess)
					ancienMess0 = mess
					break
				}
			}	
		case mess := <-clientsPresents[1].receiveChannel:
			for i:=1;i<=8;i++ {
				if string(mess) == "4"+strconv.Itoa(1)+strconv.Itoa(i) && mess != ancienMess1 {
					writeToClients(clientsPresents,mess)
					ancienMess1 = mess
					break
				}
			}
		case mess := <-clientsPresents[2].receiveChannel:
			for i:=1;i<=8;i++ {
				if string(mess) == "4"+strconv.Itoa(2)+strconv.Itoa(i) && mess != ancienMess2 {
					writeToClients(clientsPresents,mess)
					ancienMess1 = mess
					break
				}
			}
		case mess := <-clientsPresents[3].receiveChannel:
			for i:=1;i<=8;i++ {
				if string(mess) == "4"+strconv.Itoa(3)+strconv.Itoa(i) && mess != ancienMess3 {
					writeToClients(clientsPresents,mess)
					ancienMess1 = mess
					break
				}
			}
		default:
			break
		}	
	}
	
}

func chooseRunner() {
		
	runnerschose := [4]bool{}
	runnersColor := ""

	go playerSelector()

	ancienMess0 := ""
	ancienMess1 := ""
	ancienMess2 := ""
	ancienMess3 := ""
	
	for {

		select{
		case mess := <-clientsPresents[0].receiveChannel:
			if string(mess)[:2] == "3"+strconv.Itoa(0) && mess != ancienMess0{
				runnerschose[0] = true
				runnersColor += string(string(mess)[2])
				ancienMess0 = mess
			}
		case mess := <-clientsPresents[1].receiveChannel:
			if string(mess)[:2] == "3"+strconv.Itoa(1) && mess != ancienMess1 {
				runnerschose[1] = true
				runnersColor += string(string(mess)[2])
				ancienMess1 = mess
			}
		case mess := <-clientsPresents[2].receiveChannel:
			if string(mess)[:2] == "3"+strconv.Itoa(2) && mess != ancienMess2 {
				runnerschose[2] = true
				runnersColor += string(string(mess)[2])
				ancienMess2 = mess
			}
		case mess := <-clientsPresents[3].receiveChannel:
			if string(mess)[:2] == "3"+strconv.Itoa(3) && mess != ancienMess3 {
				runnerschose[3] = true
				runnersColor += string(string(mess)[2])
				ancienMess3 = mess
			}	
		}

		c := 0
		for i := range runnerschose {
			if runnerschose[i] == true {
				c++
			}
		}

		log.Println(runnerschose)
		log.Println(c)
		log.Println(runnersColor)
		if c == len(runnerschose) {
			writeToClients(clientsPresents,"400"+runnersColor)
			return
		}
	}
} 

func checkPos() {

	ancienMess0 := ""
	ancienMess1 := ""
	ancienMess2 := ""
	ancienMess3 := ""

	select {
	case mess := <-clientsPresents[0].receiveChannel:
		if mess != ancienMess0 && string(mess[:2]) == "51" {
			for i := range clientsPresents {
				if i != 0 {
					writeMessage(clientsPresents[i],"9"+mess[2:])
				}
			}
		}
		ancienMess0 = mess
	case mess := <-clientsPresents[1].receiveChannel:
		if mess != ancienMess1 && string(mess[:2]) == "51" {
			for i := range clientsPresents {
				if i != 1 {
					writeMessage(clientsPresents[i],"9"+mess[2:])
				}
			}
		}
		ancienMess1 = mess
	case mess := <-clientsPresents[2].receiveChannel:
		if mess != ancienMess2 && string(mess[:2]) == "51" {
			for i := range clientsPresents {
				if i != 2 {
					writeMessage(clientsPresents[i],"9"+mess[2:])
				}
			}
		}
		ancienMess2 = mess
	case mess := <-clientsPresents[3].receiveChannel:
		if mess != ancienMess3 && string(mess[:2]) == "51" {
			for i := range clientsPresents {
				if i != 3 {
					writeMessage(clientsPresents[i],"9"+mess[2:])
				}
			}
		}
		ancienMess3 = mess
	default:
		break
	}
}

func checkArrival() {

	ClientsFinished := [4]bool{}

	for {

		// checkPos()

		for i,client := range clientsPresents {	
			if string(<-client.receiveChannel) == "50"+strconv.Itoa(i) {	
				ClientsFinished[i] = true
				client.good = true
			}
		}

		log.Println(ClientsFinished)
		
		c := 0
		for i := range ClientsFinished {
			if ClientsFinished[i] == true || clientsPresents[i].good {
				c++
			}
		}

		log.Println("c2 = "+strconv.Itoa(c))
		if c == len(ClientsFinished) {
			writeToClients(clientsPresents,"600")
			break
		}
	}
}

func handleResults() {
	ClientsWantToRestart := [4]bool{}

	emptyChannel()

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


func emptyChannel() {
	for _,client := range clientsPresents {
		for len(client.receiveChannel) > 0 {
			<-client.receiveChannel
		}
	}
}