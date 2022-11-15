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

func (g *Game)connexion(ip string) (client net.Conn) {
	conn, err := net.Dial("tcp", ip + ":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}

	go g.readFromServer()

	return conn
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
	}
}


func (g *Game)HandleWelcomeScreenMulti() (bool) {
	if g.conn == nil && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.conn = g.connexion(ip_reseau)
	}

	select {
	case message := <- g.receiveChannel:
		if message == "200" {
			return true
		}
	default:
	}
	return false
}