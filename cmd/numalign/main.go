package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fromanirh/numalign/internal/pkg/numalign"
)

func main() {
	var err error
	hours := 24 // default
	val := os.Getenv("NUMALIGN_SLEEP_HOURS")
	if val != "" {
		hours, err = strconv.Atoi(val)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	numalign.Execute()
	time.Sleep(time.Duration(hours) * time.Hour)
}
