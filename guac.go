package guac

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/howeyc/fsnotify"
)

// Run blocks until the process receives SIGINT or SIGTERM, allowing WatchPath
// to run.
func Run() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Printf("\nStopping watch.\n")
		done <- true
	}()

	<-done
}

// WatchPath sets up a watcher on srcDir and all child directories. fn is
// executed whenever a folder being watched emits a fs event.
func WatchPath(srcDir string, fn func() (*os.File, error)) {
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
			log.Printf("Watching for file changes in %s\n", path)
		}
		return nil
	})

	go func() {
		defer watcher.Close()
		for {
			select {
			case <-watcher.Event:
				if _, err := fn(); err != nil {
					log.Println("error:", err)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
}
