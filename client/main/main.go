package main

import (
	"flag"
	"github.com/havr/lottery/client"
	"github.com/havr/lottery/game"
	"github.com/havr/lottery/util"
	"log"
	"sync"
)

func player(server string, initial int, fee int, running *sync.WaitGroup, interrupt chan struct{}) {
	defer running.Done()
	money := initial
	id := client.RandomId()
	isGameFree := false
	session := client.NewSession(server, id)
	for {
		select {
		case <-interrupt:
			return
		default:
		}
		var delta int
		if isGameFree {
			delta = 0
			isGameFree = false
		} else {
			delta = fee
		}
		money -= delta
		if money < 0 {
			break
		}
		result, err := session.Play(delta)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		switch r := result.(type) {
		case game.NoWin:
		case game.Win:
			log.Println(id, "won", r.Amount)
			money += r.Amount
		case game.FreeGame:
			log.Println(id, "got free game")
			isGameFree = true
		default:
			log.Printf("unknown result type: %T\n", r)
		}
	}
}

func main() {
	util.Init()
	serv := flag.String("s", "localhost:8888", "game server")
	nclient := flag.Int("n", 100, "number of clients")
	initmoney := flag.Int("m", 10000, "initial amount of money")
	fee := flag.Int("f", 1, "game fee")
	flag.Parse()

	var running sync.WaitGroup
	interrupt := make(chan struct{}, 1)
	for i := 0; i < *nclient; i++ {
		running.Add(1)
		go player(*serv, *initmoney, *fee, &running, interrupt)
	}

	util.WaitForInterrupt()
	close(interrupt)
	running.Wait()
}
