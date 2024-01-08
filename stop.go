package main

import (
	"encoding/json"
	"os"
)

// define a structure to hold the stopwords
type Stopwords struct {
	Words []string `json:"stopwords"`
}

func LoadStopwords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stopwords Stopwords
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&stopwords); err != nil {
		return nil, err
	}

	return stopwords.Words, nil
}

func IsStopword(word string, stopwords []string) bool {
	for _, stopword := range stopwords {
		if word == stopword {
			return true // word is a stopword
		}
	}
	return false // word is not a stopword
}
