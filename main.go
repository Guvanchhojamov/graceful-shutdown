package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{
		Addr:    ":8181",
		Handler: nil,
	}
	var t int
	signals := make(chan os.Signal, 1)
	StartServerWithGracefulShutDown(context.Background(), server, signals)
	for {
		if t != 10 {
			t += 1
			time.Sleep(time.Millisecond * 500)
			fmt.Println("server doing a job...", t)
		} else {
			signals <- syscall.SIGINT
			close(signals)
			return
		}
	}
}

func StartServerWithGracefulShutDown(ctx context.Context, server *http.Server, signals chan os.Signal) {
	go func() {
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		<-signals
		if err := server.Shutdown(ctx); err != nil {
			log.Println("shutdown error", err)
		}
		log.Println("server stopped. signal: ", signals)
		close(signals)
	}()
	fmt.Println("8181 port,server running...")
	_, err := net.Listen("tcp", ":8181")
	if err != nil {
		log.Fatal(err)
	}
}
