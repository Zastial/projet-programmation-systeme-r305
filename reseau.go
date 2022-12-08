package main

import (
	"log"
	"net"
	"bufio"
	"strings"
	"strconv"
	"time"
	"fmt"

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
	for i:=0;i<2;i++ {
		_, err := writer.WriteString(message+"\n")
		writer.Flush()
		if err!=nil{
			return
		}
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

		if message == "200" {
			return true
		}
	default:
		break
	}
	return false
}

func (g *Game) runnersColor() {

	id := strconv.Itoa(g.id_runner)

	select {
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(0)+strconv.Itoa(couleur) && id != "0" {
				g.runners[0].set_colorScheme(couleur-1)
			}
		}
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(1)+strconv.Itoa(couleur) && id != "1" {
				g.runners[1].set_colorScheme(couleur-1)
			}
		}
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(2)+strconv.Itoa(couleur) && id != "2" {
				g.runners[2].set_colorScheme(couleur-1)
			}
		}
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(3)+strconv.Itoa(couleur) && id != "3" {
				g.runners[3].set_colorScheme(couleur-1)
			}
		}
	default:
		return
	}
	
}

func (g *Game) ChooseRunnersMulti() (bool) {

	// go g.runnersColor()

	id := strconv.Itoa(g.id_runner)

	// if ebiten.IsKeyPressed(ebiten.KeyRight) && !g.good {
	// 	colorSchemeCurrentlyOn := strconv.Itoa(g.runners[g.id_runner].get_colorScheme()+1)
	// 	g.writeToServer("4"+id+colorSchemeCurrentlyOn)	
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyLeft) && !g.good {
	// 	colorSchemeCurrentlyOn := strconv.Itoa(g.runners[g.id_runner].get_colorScheme()+1)
	// 	g.writeToServer("4"+id+colorSchemeCurrentlyOn)	
	// }

	if (g.runners[g.id_runner].ManualChoose() && !g.good) {
		couleur := strconv.Itoa(g.runners[g.id_runner].get_colorScheme())
		g.writeToServer("3"+id+couleur)
		g.good = true
	}
	
	select {
	case mess := <-g.receiveChannel:
		log.Println("Message reçu : ", string(mess))
		if mess[:3] == "400" {
			for i := range g.runners {
				if i != g.id_runner {
					id, _ := strconv.Atoi(string(mess[3+i]))
					g.runners[i].set_colorScheme(id)
				}
			}
			g.good = false
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
	}
}


func (g *Game) checkPosPlayers() {

	// for i := range g.runners {
	// 	if i != g.id_runner {
	// 		g.runners[i].UpdateSpeed(false)
	// 	}
	// }

	fmt.Println("oguidfg")

	select {
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "90" && g.id_runner != 0 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[0].set_speed(playerSpeed)
				g.runners[0].UpdateSpeed(true)
				g.runners[0].UpdatePos()
			}
		}
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "91" && g.id_runner != 1 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[1].set_speed(playerSpeed)
				g.runners[1].UpdateSpeed(true)
				g.runners[1].UpdatePos()
			}
		}
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "92" && g.id_runner != 2 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[2].set_speed(playerSpeed)
				g.runners[2].UpdateSpeed(true)
				g.runners[2].UpdatePos()
			}
		}
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "93" && g.id_runner != 3 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[3].set_speed(playerSpeed)
				g.runners[3].UpdateSpeed(true)
				g.runners[3].UpdatePos()
			}
		}
	default:
		break
	}
	
}

func (g *Game) CheckArrivalMulti() (finished bool) {

	// g.checkPosPlayers()

	finished = g.runners[g.id_runner].arrived
	id := strconv.Itoa(g.id_runner)
	// speed := g.runners[g.id_runner].get_speed()

	for i := range g.runners {
		g.runners[i].CheckArrival(&g.f)
	}

	// if (!finished && inpututil.IsKeyJustPressed(ebiten.KeySpace)) {
	// 	s := fmt.Sprintf("%f", speed)
	// 	g.writeToServer("51"+id+s)
	// }

	finished = g.runners[g.id_runner].arrived
	if finished && !g.good {
		g.writeToServer("50"+id)
		g.good = true
	}

	select {
	case mess := <-g.receiveChannel:
		if mess == "600" {
			g.good = false
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