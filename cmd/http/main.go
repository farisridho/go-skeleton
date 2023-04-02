package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/farisridho/go-skeleton/interface/http"
)

func main() {
	fmt.Printf("HTTP API: go-skeleton v1")
	initGracefullShutdown()
	err := http.StartHttpServer()
	if err != nil {
		panic(err)
	}
}

func initGracefullShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)

	go func() {
		<-s
		fmt.Println("Shutting Down Gracefully......")
		//clean up
		os.Exit(1)
	}()
}
