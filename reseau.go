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

/*
	Cette fonction gère la connexion à un serveur sur l'adresse IP spécifiée dans le champ IP de l'objet Game.
	La connexion est établie en utilisant le protocole TCP sur le port 8080.
	Si la connexion réussit, l'objet conn est enregistré dans le champ conn de l'objet Game et le champ good est mis à false.
	Enfin, elle lance une goroutine appelée readFromServer pour lire les données en provenance du serveur.
*/
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

/*
	Cette fonction envoie un message au serveur auquel l'objet Game est connecté.
	Le message est envoyé en utilisant le protocole TCP.
	Le message est envoyé deux fois pour être sûr qu'il arrive correctement au serveur.
*/

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

/*
	Cette fonction lit les données en provenance du serveur auquel l'objet Game est connecté.
	Les données sont lues en continu dans une boucle infinie jusqu'à ce que la connexion soit fermée.
	Si une erreur survient pendant la lecture, elle est enregistrée dans les journaux.
	Chaque message lu est envoyé au canal receiveChannel de l'objet Game.
*/
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

/*
	Cette fonction gère l'écran de bienvenue en mode multijoueur.
	Si la touche ESPACE est pressée et que l'objet Game n'est pas connecté à un serveur, elle lance la fonction connexion pour établir la connexion.
	Ensuite, elle utilise une instruction select pour attendre les messages provenant du serveur via le canal receiveChannel de l'objet Game.
	Si un message est reçu, la fonction vérifie si le message est un nombre entier compris entre 1 et 4 inclus,
	auquel cas elle enregistre l'ID du joueur dans le champ id_runner de l'objet Game et le nombre de joueurs dans le champ nbRunner.

	Si le message est "200", la fonction renvoie true.

*/
func (g *Game)HandleWelcomeScreenMulti() (bool) {
	if g.conn == nil && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.connexion()
	}

	select {
	case message := <- g.receiveChannel:

		id, _:= strconv.Atoi(message)
		if id <=4 && id >= 1 {
			g.id_runner = id-1
			g.nbRunner = id
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
	colorSchemeCurrentlyOn := strconv.Itoa(g.runners[g.id_runner].get_colorScheme()+1)

	ancienMess0 := ""
	ancienMess1 := ""

	if ebiten.IsKeyPressed(ebiten.KeyRight) && !g.good && colorSchemeCurrentlyOn != ancienMess0  {
		g.writeToServer("4"+id+colorSchemeCurrentlyOn)
		ancienMess0 = "4"+id+colorSchemeCurrentlyOn
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && !g.good && colorSchemeCurrentlyOn != ancienMess1 {
		g.writeToServer("4"+id+colorSchemeCurrentlyOn)
		ancienMess1 = colorSchemeCurrentlyOn
	}

	select {
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(0)+strconv.Itoa(couleur) {
				g.runners[0].set_colorScheme(couleur-1)
			}
		}
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(1)+strconv.Itoa(couleur) {
				g.runners[1].set_colorScheme(couleur-1)
			}
		}
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(2)+strconv.Itoa(couleur) {
				g.runners[2].set_colorScheme(couleur-1)
			}
		}
	case mess := <-g.receiveChannel:
		for couleur:=1;couleur<=8;couleur++ {
			if string(mess) == "4"+strconv.Itoa(3)+strconv.Itoa(couleur) {
				g.runners[3].set_colorScheme(couleur-1)
			}
		}
	default:
		return
	}
	
}

/*
	Cette fonction gère l'écran du choix des personnages.
	Si le choix du joueur est fait, le client envoie un message au serveur pour lui envoyer la couleur définie.
	Tant que le client n'a pas reçu la confirmation du serveur que tout le monde a choisi sa couleur, il attend.
	Si le serveur envoie "400" alors on passe à l'écran suivant.
*/
func (g *Game) ChooseRunnersMulti() (bool) {

	id := strconv.Itoa(g.id_runner)

	if !g.good {
		g.runnersColor()
	}

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

/*
	Met à jour le joueur sur l'écran
*/
func (g *Game) UpdateRunnersMulti() {
	g.runners[g.id_runner].ManualUpdate()
}


func (g *Game) checkPosPlayers() {

	for i := range g.runners {
		if i != g.id_runner {
			g.runners[i].UpdateSpeed(false)
			g.runners[i].UpdatePos()
		}
	}

	select {
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "90" && g.id_runner != 0 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[0].set_speed(playerSpeed)
				g.runners[0].UpdateSpeed(true)
				for i:= 0; i < 30; i++ {
					g.runners[0].UpdatePos()
				}	
			}
		}
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "91" && g.id_runner != 1 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[1].set_speed(playerSpeed)
				g.runners[1].UpdateSpeed(true)
				for i:= 0; i < 30; i++ {
					g.runners[1].UpdatePos()
				}	
			}
		}
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "92" && g.id_runner != 2 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[2].set_speed(playerSpeed)
				g.runners[2].UpdateSpeed(true)
				for i:= 0; i < 30; i++ {
					g.runners[2].UpdatePos()
				}	
			}
		}
	case msg := <-g.receiveChannel :
		if string(msg[:2]) == "93" && g.id_runner != 3 {
			playerSpeed,_ := strconv.ParseFloat(string(msg[2:]), 64)
			if string(msg) != "" || string(msg) != " " {
				g.runners[3].set_speed(playerSpeed)
				g.runners[3].UpdateSpeed(true)
				for i:= 0; i < 30; i++ {
					g.runners[3].UpdatePos()
				}	
			}
		}
	default:
		break
	}
}

/*
	Cette fonction gère la course entre les joueurs.
	Si le client n'est pas arrivé, il envoie un message au serveur chaque fois que la touche ESPACE est pressée afin d'envoyer sa vitesse au serveur,
	et dès qu'il est arrivé, il envoie '500' et attend.
	Tant que le client n'a pas reçu la confirmation du serveur que tout le monde est arrivé, il continue d'attendre.
	Si le serveur envoie "600" alors on passe à l'écran suivant.
*/
func (g *Game) CheckArrivalMulti() (finished bool) {

	if !g.good {
		g.checkPosPlayers()
	}

	finished = g.runners[g.id_runner].arrived
	id := strconv.Itoa(g.id_runner)
	speed := g.runners[g.id_runner].get_speed()

	for i := range g.runners {
		g.runners[i].CheckArrival(&g.f)
	}

	if (!finished && inpututil.IsKeyJustPressed(ebiten.KeySpace)) {
		s := fmt.Sprintf("%f", speed)
		g.writeToServer("51"+id+s)
	}

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

/*
	Cette fonction gère l'affichage des résultats et le redémarrage du jeu.
	Si le client n'est pas arrivé, il envoie un message au serveur chaque fois que la touche ESPACE est pressée afin d'envoyer sa vitesse au serveur,
	et dès qu'il est arrivé, il envoie '500' et attend.
	Tant que le client n'a pas reçu la confirmation du serveur que tout le monde est arrivé, il continue d'attendre.
	Si le serveur envoie "600" alors on passe à l'écran suivant.
*/
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




