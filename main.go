package main

import (
	"fmt"
	"log"
	"os"
	"talker/config"
)

var (
	GitTag string = "dev"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
