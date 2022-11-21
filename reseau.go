package main

import (
	"log"
	"net"
	"bufio"
	"strings"
	// "time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var ip_reseau string

func (g *Game)connexion() {
	conn, err := net.Dial("tcp", g.IP + ":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	g.conn = conn
	go g.readFromServer()
}

func (g *Game)writeToServer(message string) {
    writer := bufio.NewWriter(g.conn)
    _, err := writer.WriteString(message+"\n")
    writer.Flush()
    if err!=nil{
        return
    }
}

func (g *Game)readFromServer() (text string){
	reader := bufio.NewReader(g.conn)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Println("read error:", err)
		}
		strip := strings.TrimSuffix(s, "\n")

		log.Println("received message from server : ", strip)

		g.receiveChannel <- strip
		log.Println("Le message enregistré est : ", <-g.receiveChannel)

		return strip
	}
}


func (g *Game)HandleWelcomeScreenMulti() (bool) {
	if g.conn == nil && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		go g.connexion()
	}

	select {
	case message := <- g.receiveChannel:
		log.Println("Waiting for the message..")
		if message == "200" {
			return true
		}
	default:
		break
	}
	return false
}