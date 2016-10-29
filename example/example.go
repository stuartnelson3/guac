package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stuartnelson3/guac"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher, err := guac.NewWatcher(ctx, "./concat", func() error {
		fmt.Println("change detected")
		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	go watcher.Run()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("watching")
	<-done
	fmt.Println("watch has ended")
}
