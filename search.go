package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/kljensen/snowball"
)

func DbSearchHandler(w http.ResponseWriter, r *http.Request, index *dbSql) {

	r.ParseForm()

	var isBigram = false
	isWild := false

	termValues, ok := r.Form["term"]

	if !ok || len(termValues) == 0 {
		fmt.Fprintf(w, "No search term entered")
		return
	}

	images, okimg := r.Form["images"]

	if okimg || len(images) > 0 {

		imgs := index.SearchImg(termValues[0])
		tmpl, err := template.ParseFiles("images.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, imgs)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		str := termValues[0]

		wildcard, hm := r.Form["wildcard"]

		if hm || len(wildcard) > 1 {
			isWild = true
		}

		if strings.HasPrefix(termValues[0], "*") {
			isWild = true
			str = strings.TrimPrefix(str, "*")
		}

		var bigram = []string{}

		bigram = strings.Split(str, " ")

		term1 := ""
		term2 := ""
		if len(bigram) > 1 {
			isBigram = true
			term1 = bigram[0]
			term2 = bigram[1]
		} else {
			term1 = str
		}

		term, err := snowball.Stem(term1, "english", true)
		if err != nil {
			fmt.Println(err)
		}
		term2, err = snowball.Stem(term2, "english", true)
		if err != nil {
			fmt.Println(err)
		}

		tfidf := DbCalculateTFIDF(term, index, isBigram, term2, isWild)

		tmpl, err := template.ParseFiles("temp.html")
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(w, tfidf)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// starts the web server
func DbSearch(idx *dbSql) {
	// create a handler for serving static files (e.g., HTML, CSS) from the "static" directory
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer) // Handle requests at the root URL ("/")

	// create a handler for the "/search" endpoint
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		DbSearchHandler(w, r, idx)
	})

	// start the web server on port 8080
	port := 8080
	fmt.Printf("Starting server at port %d\n", port)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
