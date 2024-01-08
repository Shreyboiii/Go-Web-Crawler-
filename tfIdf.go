package main

import (
	"sort"
)

type result struct {
	Term  string
	Title string
	Url   string
	TFIDF float64
}

type tempMap struct {
	Title string
	TFIDF float64
}

type sortResults []result

func (r sortResults) Len() int           { return len(r) }
func (r sortResults) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r sortResults) Less(i, j int) bool { return r[i].TFIDF > r[j].TFIDF }

func DbCalculateTFIDF(term string, idx *dbSql, isBigram bool, term2 string, isWild bool) []result {
	temp := map[string]tempMap{}

	if isBigram {

		biFreq := idx.SearchBi(term, term2)
		if isWild {
			biFreq = idx.SearchBiWild(term, term2)
		}

		for url, data := range biFreq {

			tc := float64(idx.WordsInDoc(url))
			TF := float64(data.count) / tc
			df := float64(len(biFreq)) / (float64(idx.DocsInCorp()))
			IDF := 1.0 / df
			TFIDF := TF * IDF
			temp[url] = tempMap{
				data.title,
				TFIDF,
			}
			// fmt.Printf("url: %s word count: %d, ", url, idx.WordsInDoc(url))
			// fmt.Printf("term count: %d, ", termCount)
			// fmt.Printf("docs count: %d, score: %g\n", idx.DocsInCorp(), TFIDF)
		}

		var results2 []result
		newTerm := term + " " + term2
		for key, value := range temp {
			results2 = append(results2, result{newTerm, value.Title, key, value.TFIDF})
		}
		return results2
	} else {

		freq := idx.Search(term)
		if isWild {
			freq = idx.SearchWild(term)
		}

		for url, data := range freq {

			tc := float64(idx.WordsInDoc(url))
			TF := float64(data.count) / tc
			df := float64(len(freq)) / (float64(idx.DocsInCorp()))
			IDF := 1.0 / df
			TFIDF := TF * IDF
			temp[url] = tempMap{
				data.title,
				TFIDF,
			}
		}

		var results []result

		for key, value := range temp {
			results = append(results, result{term, value.Title, key, value.TFIDF})
		}

		sort.Sort(sortResults(results))
		return results
	}

}
