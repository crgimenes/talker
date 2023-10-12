package main

import (
	"fmt"
	"log"
	"os"
	"talker/config"
	"talker/session"
	"time"
)

const cookieName = "session_talker"

var (
	GitTag string = "dev"
	sc     *session.Control
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	session.Create(cookieName)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			session.SC.RemoveExpired()
		}
	}()

}
