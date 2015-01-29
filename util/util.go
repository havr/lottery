package util

import (
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func Init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().Unix())
}

func WaitForInterrupt() {
	interrupt := make(chan os.Signal, 1)
	defer close(interrupt)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
}
