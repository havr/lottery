package game

import (
	"math/rand"
	"sync"
	"time"
)

const (
	pairDiscardTime = 10 * time.Second
)

type Pair [2]byte

func randomByte() byte {
	return byte(rand.Int() % 256)
}

func RandomPair() Pair {
	return Pair{randomByte(), randomByte()}
}

type luckyStack struct {
	pairs []Pair
	head  int
}

func newLuckyStack(size int) *luckyStack {
	lucky := &luckyStack{
		pairs: make([]Pair, size),
	}
	for i := 0; i < size; i++ {
		lucky.pairs[i] = RandomPair()
	}
	return lucky
}

func (stack *luckyStack) Pop() Pair {
	pair := stack.pairs[stack.head]
	stack.pairs[stack.head] = RandomPair()
	stack.head = (stack.head + 1) % len(stack.pairs)
	return pair
}

type Id int64

type bid struct {
	id     Id
	pair   Pair
	fee    int
	result chan interface{}
}

type NoWin struct{}

type Win struct{ Amount int }

type FreeGame struct{}

type Game struct {
	bids      chan bid
	running   sync.WaitGroup
	interrupt chan struct{}
}

func Spawn() *Game {
	g := &Game{
		bids:      make(chan bid, 1),
		interrupt: make(chan struct{}, 1),
	}
	go g.run()
	return g
}

func (game *Game) run() {
	defer game.running.Done()

	jackpot := 0
	freeGameIds := make(map[Id]struct{})
	luckyStack := newLuckyStack(100)
	timer := time.NewTimer(pairDiscardTime)
	defer timer.Stop()

	for {
		select {
		case <-game.interrupt:
			return
		case <-timer.C:
			luckyStack.Pop()
			timer.Reset(pairDiscardTime)
		case bid := <-game.bids:
			go timer.Reset(pairDiscardTime)
			if _, isGameFree := freeGameIds[bid.id]; isGameFree {
				delete(freeGameIds, bid.id)
			} else if bid.fee == 0 {
				bid.result <- NoWin{}
				continue
			}
			jackpotWasEmpty := jackpot == 0
			jackpot += bid.fee
			lucky := luckyStack.Pop()
			if bid.pair == lucky {
				if jackpotWasEmpty {
					bid.result <- FreeGame{}
				} else {
					bid.result <- Win{Amount: jackpot}
					jackpot = 0
				}
			} else {
				bid.result <- NoWin{}
			}
		}
	}
}

func (game *Game) Bid(id Id, pair Pair, fee int) interface{} {
	result := make(chan interface{}, 1)
	defer close(result)
	game.bids <- bid{id: id, pair: pair, fee: fee, result: result}
	return <-result
}

func (game *Game) Shutdown() {
	close(game.interrupt)
	game.running.Wait()
}
