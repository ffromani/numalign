package main

import (
	"time"

	"github.com/fromanirh/numalign/cmd"
)

func main() {
	cmd.Execute()
	time.Sleep(24 * time.Hour)
}
