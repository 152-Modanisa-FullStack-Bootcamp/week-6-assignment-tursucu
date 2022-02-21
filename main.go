package main

import (
	"log"
	"wallet/server"
)

func main() {
	serve := server.NewServer()
	err := serve.StartServer(3000)
	if err != nil {
		log.Fatalln(err)
	}
}
