package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Scan a root folder and return a list of Movies (.nfo file).
func Scan(rootPath string) []*Movie {
	list := make([]*Movie, 0)

	// Walk
	err := filepath.Walk(rootPath, func(objPath string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".nfo") {
			m := scanNfo(rootPath, objPath)
			list = append(list, m)
		}
		return nil
	})

	// error or return
	if err != nil {
		panic(err)
	}
	return list
}

// scanNfo scan a nfo and return a Movie struct.
func scanNfo(rootPath, nfoPath string) *Movie {
	ret := new(Movie)

	// parse xml
	b, err := ioutil.ReadFile(nfoPath)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(b, &ret)
	if err != nil {
		panic(err)
	}

	// set paths
	ret.NfoPath = nfoPath
	ret.DirPath = filepath.Dir(ret.NfoPath)
	ret.PosterPath = filepath.Join(ret.DirPath, "poster.jpg")

	// rel dir path
	relDir, err := filepath.Rel(rootPath, ret.DirPath)
	if err != nil {
		panic(err)
	}
	ret.DirPath = relDir

	// fin
	return ret
}
