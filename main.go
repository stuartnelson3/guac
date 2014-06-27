package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	concat(*dst, *srcDir)
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
