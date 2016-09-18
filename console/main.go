package console

import (
	"log"
	"os"
	"os/signal"
)

func WaitForKill(done chan<- bool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	for _ = range ch {
		log.Println("Got `CTRL-C` from console, notify done chan")
		done <- true
		break
	}
}
