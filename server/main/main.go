package main

import (
	"flag"
	"github.com/havr/lottery/game"
	"github.com/havr/lottery/server"
	"github.com/havr/lottery/util"
	"log"
)

func main() {
	util.Init()
	host := flag.String("host", ":8888", "server hostname")
	flag.Parse()

	game := game.Spawn()
	serv, err := server.Spawn(*host, game)
	if err != nil {
		log.Println(err.Error())
		return
	}

	util.WaitForInterrupt()
	serv.Shutdown()
	game.Shutdown()
}
