package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".jpg") {
			return nil
		}

		h := sha256.New()
		h.Write([]byte(path))
		b := h.Sum(nil)
		newname := fmt.Sprintf("%x.jpg", b[0:10])

		return os.Rename(path, newname)
	})
	if err != nil {
		fmt.Println(err)
	}
}
