package concat

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Concat walks srcDir and appends all files found with extension matching ext
// to dst. Files are appended in lexical order.
func Concat(dst, srcDir, ext string) (*os.File, error) {
	f, err := os.Create(dst)
	filepath.Walk(srcDir, find(ext, dst))
	return f, err
}

func find(ext, dst string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			log.Printf("Appending %s\n", path)
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
