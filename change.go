package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func changeHash(rootPath string) string {
	h := sha256.New()

	if err := filepath.Walk(rootPath, func(objPath string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".nfo") {
			// vars
			name := f.Name()
			size := fmt.Sprintf("%d", f.Size())
			time := fmt.Sprintf("%d", f.ModTime().Unix())
			// hash
			h.Write([]byte(name))
			h.Write([]byte(","))
			h.Write([]byte(size))
			h.Write([]byte(","))
			h.Write([]byte(time))
			h.Write([]byte("|"))
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return fmt.Sprintf("<!-- %x -->", h.Sum(nil))
}
