package client

import (
	"encoding/gob"
	"github.com/havr/lottery/game"
	"github.com/havr/lottery/server"
	"math/rand"
	"net"
)

type Session struct {
	server string
	id     game.Id
}

func RandomId() game.Id {
	return game.Id(rand.Int63())
}

func NewSession(server string, id game.Id) *Session {
	return &Session{
		server: server,
		id:     id,
	}
}

func (session *Session) Play(fee int) (interface{}, error) {
	conn, err := net.Dial("tcp", session.server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	request := server.Request{Id: session.id, Pair: game.RandomPair(), Fee: fee}
	err = gob.NewEncoder(conn).Encode(&request)
	if err != nil {
		return nil, err
	}

	var response interface{}
	err = gob.NewDecoder(conn).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
