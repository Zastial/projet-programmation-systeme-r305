/*
// Implementation of a main function setting a few characteristics of
// the game window, creating a game, and launching it
*/

package main

import (
	"flag"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800 // Width of the game window (in pixels)
	screenHeight = 160 // Height of the game window (in pixels)
)

func main() {

	var getTPS bool
	flag.BoolVar(&getTPS, "tps", false, "Afficher le nombre d'appel à Update par seconde")
	flag.Parse()

	var ip string
	flag.StringVar(&ip, "serverip", "localhost", "IP du serveur")
	flag.Parse()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("BUT2 année 2022-2023, R3.05 Programmation système")

	g := InitGame()
	g.getTPS = getTPS

	connexion(ip)

	err := ebiten.RunGame(&g)
	log.Print(err)

}
