package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func main() {

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".jpg") {
			return nil
		}

		path = norm.NFC.String(path)

		h := sha256.New()
		h.Write([]byte(strings.Replace(path, ":", "/", -1)))
		b := h.Sum(nil)
		newname := fmt.Sprintf("%x.jpg", b[0:10])

		return os.Rename(path, newname)
	})
	if err != nil {
		fmt.Println(err)
	}
}
