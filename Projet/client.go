package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	//"os"
	//"strings"
)

func main() {

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("Je suis connecté")

	for {
		netData, _ := bufio.NewReader(conn).ReadString('\n')

		fmt.Print(string(netData))

		myTime := time.Now().Format(time.RFC3339) + "\n"
		conn.Write([]byte(myTime))
	}
}
