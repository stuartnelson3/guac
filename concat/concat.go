package concat

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Concat walks srcDir and appends all files found with extension ext
// to dst. Files are appended in lexical order.
func Concat(dst, srcDir, ext string) func() error {
	return func() error {
		if _, err := os.Create(dst); err != nil {
			return err
		}
		err := filepath.Walk(srcDir, find(ext, dst))
		return err
	}
}

func find(ext, dst string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			log.Printf("Appending %s to %s\n", path, dst)
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
