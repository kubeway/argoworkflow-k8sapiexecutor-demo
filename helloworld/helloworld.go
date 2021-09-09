package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	signalDebugModel := strings.ToLower(os.Getenv("SIGNAL_DEBUG_MODE"))
	switch signalDebugModel {
	case "immediately":
		log.Println("started with IMMEDIATELY")
		// wait signal
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		<-c
		return
	case "wait":
		log.Println("started with WAIT")
		// wait signal
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		s := <-c
		log.Println("signal: ", s)
		time.Sleep(time.Second * 3)
		return
	case "ignore":
		log.Println("started with ignore")
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		ticker := time.NewTicker(time.Second * 2)
		for {
			select {
			case s := <-c:
				log.Printf("ignore signal:%s\n", s)
			case tm := <-ticker.C:
				time.Sleep(time.Second)
				log.Println("The Current time is: ", tm)
			}
		}
	default:
		log.Println("started with default")
		for {
			time.Sleep(time.Second * 3)
			log.Println("The Current time is: ", time.Now())
		}
	}
}
