package main

import (
	"log"
	"net"
	"bufio"
	"strings"
	"strconv"
	"time"

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
	g.good = false
	go g.readFromServer()
}

func (g *Game)writeToServer(message string) {
    writer := bufio.NewWriter(g.conn)
    _, err := writer.WriteString(message+"\n")
    writer.Flush()
    if err!=nil{
        return
    }
	log.Println("Message envoyé au serveur : ", message)
}

func (g *Game)readFromServer() {
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
		g.connexion()
	}

	select {
	case message := <- g.receiveChannel:

		id, _:= strconv.Atoi(message)
		if id <=4 && id >= 1 {
			g.id_runner = id-1
			log.Println("You are the player : ", g.id_runner)
		}

		log.Println("Waiting for the message..")
		if message == "200" {
			return true
		}
	default:
		break
	}
	return false
}

func (g *Game) ChooseRunnersMulti() (bool) {

	if (g.runners[g.id_runner].ManualChoose()) {
		id := strconv.Itoa(g.id_runner)
		couleur := strconv.Itoa(g.runners[g.id_runner].get_colorScheme())
		g.writeToServer("3"+id+couleur)
	}

	select {
	case mess := <-g.receiveChannel:
		log.Println("Waiting for the message..")
		if mess == "400" {
			return true
		}
	default:
		break
	}
	return false
}


func (g *Game) UpdateRunnersMulti() {
	for i := range g.runners {
		if i == g.id_runner {
			g.runners[g.id_runner].ManualUpdate()
		}
		if (i != 0 && i != 1) {
			g.runners[i].RandomUpdate() //test les 2 derniers coureurs
		}
	}
}


func (g *Game) CheckArrivalMulti() (finished bool) {

	finished = false

	for i := range g.runners {

		g.runners[i].CheckArrival(&g.f)
		finished = g.runners[g.id_runner].arrived

		if finished {
			id := strconv.Itoa(g.id_runner)
			g.writeToServer("50"+id)
		}
	}

	select {
	case mess := <-g.receiveChannel:
		log.Println("Waiting for the message..")
		if mess == "600" {
			return true
		}
	default:
		break
	}
	return false
}

func (g *Game) HandleResultsMulti() bool {

	if !g.good {
		if time.Since(g.f.chrono).Milliseconds() > 1000 || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.resultStep++
			g.f.chrono = time.Now()
		}
	}

	if g.resultStep >= 4 && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep = 0
		id := strconv.Itoa(g.id_runner)
		g.writeToServer("70"+id)
		g.good = true
	}
	select {
	case mess := <-g.receiveChannel:
		log.Println("Waiting for the message..")
		if mess == "800" {
			g.good = false
			return true
		}
	default:
		break
	}
	return false
}



















// func (g *Game) getCouleurs() {

// 	message := <- g.receiveChannel

// 	words := strings.Fields(message)

// 	var couleur = ""
// 	for _, word := range words {
// 		couleur = word
// 		break
// 	}
	
// 	if string(couleur[0]) == "0" {
// 		log.Print("C'EST BON CHAKAL")
// 		g.good = true
// 	}

// 	return
// }