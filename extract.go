package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/net/html"
)

type image struct {
	alt string
	src string
}

func extract(downOut chan dOut, extractOut chan extrOut, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range downOut {

		r := bytes.NewReader(data.body)
		z, err := html.Parse(r)

		if err != nil {
			fmt.Println(err)
		}

		words := []string{}
		hrefs := []string{}
		titles := ""
		images := []image{}

		i := 0
		var enteredBody bool
		enteredBody = false

		// recursive function for traversing HTML tree
		var f func(*html.Node, *[]string, *[]string, *string, *[]image)
		// checks if letter or number
		letter := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
		f = func(n *html.Node, words *[]string, hrefs *[]string, titles *string, images *[]image) {

			nextType := n.Type // gets the type of node

			switch nextType {

			case html.ErrorNode: // error handling, eof check
				return

			case html.ElementNode:
				if n.Data == "a" {
					for _, attr := range n.Attr {
						if attr.Key == "href" { // if it is an href which is what we are looking for
							*hrefs = append(*hrefs, attr.Val) // append it to the list
						}
					}
				} else if n.Data == "img" {
					img := image{}
					for _, attr := range n.Attr {
						switch attr.Key {
						case "alt":
							img.alt = attr.Val
						case "src":
							img.src = attr.Val
						}
					}
					*images = append(*images, img)
				} else if i < 1 && n.Data == "title" {
					*titles = n.FirstChild.Data
					i++
				} else if i < 2 {
					if n.Data == "body" {
						enteredBody = true
						i++
					}
				}

			case html.TextNode: // if type is text get the text
				if enteredBody {
					word := strings.FieldsFunc(n.Data, letter)
					*words = append(*words, word...)
				}

			}
			// process child nodes
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, words, hrefs, titles, images)
			}
		}
		f(z, &words, &hrefs, &titles, &images)

		extractOut <- extrOut{
			words,
			hrefs,
			titles,
			images,
			data.body,
			data.link,
		}
	}

	close(extractOut)
}
