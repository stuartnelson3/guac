package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fazalmajid/guac"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher, err := guac.NewWatcher(ctx, "./concat", time.Second, func() error {
		fmt.Println("change detected")
		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

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
	watcher.Close()
}
