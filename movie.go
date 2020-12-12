package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"html"
	"image"
	"image/jpeg"
	"os"
)

// Movie represent a single movie (@see Scan()).
type Movie struct {
	DirPath    string
	NfoPath    string
	PosterPath string

	Plot  string `xml:"plot"`
	Added string `xml:"dateadded"`
	Title string `xml:"title"`
	Year  string `xml:"year"`

	imageCache []byte
}

// Thumbnail return a thumbnail from poster.jpg
func (m *Movie) Thumbnail(max uint) []byte {
	if m.imageCache != nil {
		return m.imageCache
	}

	// read image
	r, err := os.Open(m.PosterPath)
	if err != nil {
		return nil
	}
	defer r.Close()

	im, _, err := image.Decode(r)
	if err != nil {
		return nil
	}

	// resize
	im = resize.Thumbnail(max, max, im, resize.MitchellNetravali)

	// write new image
	w := new(bytes.Buffer)
	err = jpeg.Encode(w, im, nil)
	if err != nil {
		return nil
	}

	// success
	m.imageCache = w.Bytes()
	return m.imageCache
}

// HtmlImage return a thumbnail from poster.jpg as base64 html image.
func (m *Movie) HtmlImage(max uint) string {
	b := m.Thumbnail(max)
	b64 := base64.StdEncoding.EncodeToString(b)
	text := fmt.Sprintf("<img src='data:image/png;base64,%s' />", b64)
	return text
}

func (m *Movie) TableRow(max uint) string {
	link := fmt.Sprintf("<a href='%s?html'>%s</a>", m.DirPath, html.EscapeString(m.Title))
	return fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
		m.HtmlImage(max), link, html.EscapeString(m.Year), html.EscapeString(m.Plot), html.EscapeString(m.Added))
}

func (m *Movie) Div(max uint) string {
	img := fmt.Sprintf("<a href='%s?html'>%s</a>", m.DirPath, m.HtmlImage(max))
	title := fmt.Sprintf("<a href='%s?html'>%s</a>", m.DirPath, html.EscapeString(m.Title))
	return fmt.Sprintf("<div class='box'>%s<br>%s (%s)</div>\n", img, title, html.EscapeString(m.Year))
}
