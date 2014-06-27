package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/howeyc/fsnotify"
)

func main() {
	var (
		showHelp = flag.Bool("h", false, "show this help")
		srcDir   = flag.String("src", "", "the source directory for your js")
		dst      = flag.String("dst", "", "the file to write to")
	)
	flag.Parse()

	if *showHelp || *srcDir == "" || *dst == "" {
		flag.Usage()
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	filepath.Walk(*srcDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			err = watcher.Watch(path)
			if err != nil {
				return err
			}
			log.Printf("Watching %s.\n", path)
		}
		return nil
	})

	go func() {
		log.Printf("Watching for %s file changes in %s", ".js", *srcDir)
		for {
			select {
			case <-watcher.Event:
				concat(*dst, *srcDir)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Printf("\nStopping watch.\n")
		done <- true
	}()

	concat(*dst, *srcDir)

	<-done
}

func concat(dst, srcDir string) (*os.File, error) {
	f, err := os.Create(dst)
	filepath.Walk(srcDir, find(".js", dst))
	return f, err
}

func find(ext, dst string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			fmt.Printf("Appending %s\n", path)
			dst, err := os.OpenFile(dst, os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			defer dst.Close()

			src, err := os.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()

			io.Copy(dst, src)
		}
		return nil
	}
}
