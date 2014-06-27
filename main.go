package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/howeyc/fsnotify"
	"github.com/stuartnelson3/guac/concat"
)

func main() {
	var (
		showHelp = flag.Bool("h", false, "show this help")
		srcDir   = flag.String("src", "", "the source directory for your js")
		dst      = flag.String("dst", "", "the file to write to")
		ext      = flag.String("ext", ".js", "the file extension to update")
	)
	flag.Parse()

	if *showHelp || *srcDir == "" || *dst == "" {
		flag.Usage()
		return
	}

	WatchPath(*srcDir, *dst, *ext, concat.Concat)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Printf("\nStopping watch.\n")
		done <- true
	}()

	concat.Concat(*dst, *srcDir, *ext)

	<-done
}

// WatchPath sets up a watcher on srcDir and all child directories. fn is
// executed with arguments srcDir, a destination, and an extension whenever a
// folder being watched emits a fs event.
func WatchPath(srcDir, dst, ext string, fn func(dst, srcDir, ext string) (*os.File, error)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
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
		defer watcher.Close()
		log.Printf("Watching for %s file changes in %s", ext, srcDir)
		for {
			select {
			case <-watcher.Event:
				fn(dst, srcDir, ext)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
}
