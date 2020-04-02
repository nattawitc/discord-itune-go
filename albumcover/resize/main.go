package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func main() {

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".jpg") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error: %v   path: %v", err, path)
		}

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("error: %v   path: %v", err, path)
		}

		file.Close()

		if img.Bounds().Size().X > 512 && img.Bounds().Size().Y > 512 {
			return nil
		}

		m := resize.Resize(uint(img.Bounds().Size().X*2), 0, img, resize.Lanczos3)

		out, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("error: %v   path: %v", err, path)
		}
		defer out.Close()

		err = jpeg.Encode(out, m, nil)
		if err != nil {
			return fmt.Errorf("error: %v   path: %v", err, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
