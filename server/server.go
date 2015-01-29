package server

import (
	"encoding/gob"
	"github.com/havr/lottery/game"
	"log"
	"net"
	"sync"
)

func init() {
	gob.Register(game.NoWin{})
	gob.Register(game.Win{})
	gob.Register(game.FreeGame{})
}

type Request struct {
	Id   game.Id
	Pair game.Pair
	Fee  int
}

func (server *Server) handleConnection(conn net.Conn) {
	defer server.clients.Done()
	defer conn.Close()

	var request Request
	err := gob.NewDecoder(conn).Decode(&request)
	if err != nil {
		log.Println(err.Error())
		return
	}

	result := server.game.Bid(request.Id, request.Pair, request.Fee)

	err = gob.NewEncoder(conn).Encode(&result)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

type Server struct {
	listener  net.Listener
	game      *game.Game
	running   sync.WaitGroup
	clients   sync.WaitGroup
	interrupt chan struct{}
}

func (server *Server) waitForInterrupt() {
	<-server.interrupt
	server.listener.Close()
	server.clients.Wait()
	server.running.Done()
}

func (server *Server) run() {
	go server.waitForInterrupt()
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			break
		}
		server.clients.Add(1)
		go server.handleConnection(conn)
	}
}

func Spawn(listenTo string, game *game.Game) (*Server, error) {
	listener, err := net.Listen("tcp", listenTo)
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener:  listener,
		game:      game,
		interrupt: make(chan struct{}),
	}
	go s.run()
	return s, nil
}

func (server *Server) Shutdown() {
	close(server.interrupt)
	server.running.Wait()
}
