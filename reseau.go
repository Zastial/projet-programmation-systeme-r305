package main

import (
	"log"
	"net"
)

func connexion(ip string) {
	conn, err := net.Dial("tcp", ip + ":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()
}