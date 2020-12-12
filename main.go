package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "1.1"

const (
	maxImageSize = 200
	templateHash = "<!-- HASH -->"
	//--------
	insertDivs   = "<!-- INSERT DIVS -->"
	templateDivs = "data/template-over.htm"
	outDivs      = "index.htm"
	//--------
	insertRows   = "<!-- INSERT ROWS -->"
	templateRows = "data/template-list.htm"
	outRows      = "index2.htm"
)

// main write two files (outDivs and outRows) to the root folder.
func main() {
	// param
	if len(os.Args) < 2 {
		println("VERSION", VERSION, ": call", filepath.Base(os.Args[0]), "<root path>")
		os.Exit(1)
	}
	var root = os.Args[1]

	// check changes
	hash := changeHash(root)
	b, err := ioutil.ReadFile(filepath.Join(root, outDivs))
	if err == nil && strings.Contains(string(b), hash) {
		os.Exit(0) // no changes -> exit
	}

	// scan root
	rows := new(strings.Builder)
	divs := new(strings.Builder)
	for _, m := range Scan(root) {
		rows.WriteString(m.TableRow(maxImageSize))
		divs.WriteString(m.Div(maxImageSize))
	}

	// ----------------------------------------------

	// read template
	b, err = ioutil.ReadFile(templateRows)
	if err != nil {
		panic(err)
	}
	// replace
	html := strings.ReplaceAll(string(b), templateHash, hash)
	html = strings.ReplaceAll(html, insertRows, rows.String())
	// write file
	err = ioutil.WriteFile(filepath.Join(root, outRows), []byte(html), 0600)
	if err != nil {
		panic(err)
	}

	// ----------------------------------------------

	// read template
	b, err = ioutil.ReadFile(templateDivs)
	if err != nil {
		panic(err)
	}
	// replace
	html = strings.ReplaceAll(string(b), templateHash, hash)
	html = strings.ReplaceAll(html, insertDivs, divs.String())
	// write file
	err = ioutil.WriteFile(filepath.Join(root, outDivs), []byte(html), 0600)
	if err != nil {
		panic(err)
	}
}
