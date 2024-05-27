package main

import (
	"log"

	"github.com/hidu/http-echo-server/internal"
)

func main() {
	internal.Run()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}
