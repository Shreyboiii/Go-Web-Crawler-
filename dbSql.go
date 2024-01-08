package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/kljensen/snowball"
	_ "modernc.org/sqlite"
)

type webInfo struct {
	title string
	count int
}

type imgInfo struct {
	URL   string
	Alt   string
	Src   string
	Title string
}

type dbSql struct {
	db *sql.DB
}

var mu sync.Mutex

func NewDb() dbSql {

	db, err := sql.Open("sqlite", "proj6.db")
	if err != nil {
		log.Fatal("SQL open err ", err)
	}

	idx := dbSql{
		db: db,
	}

	if _, err = db.Exec(`
	PRAGMA foreign_keys = ON;
	`); err != nil {
		log.Fatal(err)
	}

	//creating a terms table
	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS terms(
		id integer PRIMARY KEY,
		name text UNIQUE
	);
	`); err != nil {
		log.Fatal("create terms table ", err)
	}

	//creating a urls table
	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS urls(
		id integer PRIMARY KEY,
		name text UNIQUE,
		title text,
		count integer
	);
	`); err != nil {
		log.Fatal("create urls table ", err)
	}

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS results(
		id integer PRIMARY KEY,
		term_id integer,
		url_id integer,
		term_count integer,
		FOREIGN KEY (term_id) REFERENCES terms(id),
		FOREIGN KEY (url_id) REFERENCES urls(id)
	);
	`); err != nil {
		log.Fatal("create results table ", err)
	}

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS bigrams(
		id integer PRIMARY KEY,
		term1 integer,
		term2 integer,
		url_id integer,
		freq_count integer,
		FOREIGN KEY (term1) REFERENCES terms(id),
		FOREIGN KEY (term2) REFERENCES terms(id),
		FOREIGN KEY (url_id) REFERENCES urls(id)
	);
	`); err != nil {
		log.Fatal("create bigrams table ", err)
	}

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS images(
		id integer PRIMARY KEY,
		url text,
		alt text, 
		src text,
		title text
	);
	`); err != nil {
		log.Fatal("create images table ", err)
	}

	return idx

}
func (idx dbSql) SearchBiWild(term string, term2 string) map[string]webInfo {

	var termID int64
	var term2ID int64
	var freq = map[string]webInfo{}

	err := idx.db.QueryRow(`
	SELECT id FROM terms WHERE name LIKE ?;
	`, term).Scan(&termID)

	if err != nil {
		return freq
	}

	err = idx.db.QueryRow(`
	SELECT id FROM terms WHERE name LIKE ?;
	`, term2).Scan(&term2ID)

	if err != nil {
		return freq
	}

	bi_id_rows, err := idx.db.Query(`
	SELECT id FROM bigrams WHERE term1 = ? AND term2= ?;
	`, termID, term2ID)

	if err != nil {
		return freq
	}

	for bi_id_rows.Next() {
		var bId int64
		err := bi_id_rows.Scan(&bId)
		if err != nil {
			log.Fatal("searchBi bi_id_rows.Scan error ", err)
		}
		newRow := idx.db.QueryRow(`
			SELECT urls.name, urls.title, bigrams.freq_count FROM urls
			INNER JOIN bigrams ON urls.id = bigrams.url_id AND bigrams.id = ?;`, bId)
		var url string
		var title string
		var count int
		err = newRow.Scan(&url, &title, &count)
		if err != nil {
			log.Fatal("searchBi newRow.scan error ", err)
		}
		freq[url] = webInfo{
			title,
			count,
		}
	}

	return freq
}

func (idx dbSql) SearchBi(term string, term2 string) map[string]webInfo {

	var termID int64
	var term2ID int64
	var freq = map[string]webInfo{}

	err := idx.db.QueryRow(`
	SELECT id FROM terms WHERE name = ?;
	`, term).Scan(&termID)

	if err != nil {
		return freq
	}

	err = idx.db.QueryRow(`
	SELECT id FROM terms WHERE name = ?;
	`, term2).Scan(&term2ID)

	if err != nil {
		return freq
	}

	bi_id_rows, err := idx.db.Query(`
	SELECT id FROM bigrams WHERE term1 = ? AND term2= ?;
	`, termID, term2ID)

	if err != nil {
		return freq
	}

	for bi_id_rows.Next() {
		var bId int64
		err := bi_id_rows.Scan(&bId)
		if err != nil {
			log.Fatal("searchBi bi_id_rows.Scan error ", err)
		}
		newRow := idx.db.QueryRow(`
			SELECT urls.name, urls.title, bigrams.freq_count FROM urls
			INNER JOIN bigrams ON urls.id = bigrams.url_id AND bigrams.id = ?;`, bId)
		var url string
		var title string
		var count int
		err = newRow.Scan(&url, &title, &count)
		if err != nil {
			log.Fatal("searchBi newRow.scan error ", err)
		}
		freq[url] = webInfo{
			title,
			count,
		}
	}

	return freq
}

func (idx dbSql) SearchWild(term string) map[string]webInfo {

	var termID int64
	var freq = map[string]webInfo{}

	err := idx.db.QueryRow(`
	SELECT id FROM terms WHERE name LIKE ?;
	`, term).Scan(&termID)

	if err != nil {
		return freq
	}

	rows, err := idx.db.Query(`
	SELECT urls.name, urls.title, results.term_count FROM urls
	INNER JOIN results ON urls.id = results.url_id AND results.term_id = ?;`, termID)

	if err != nil {
		log.Fatal("searchWild query error ", err)
	}

	for rows.Next() {
		var url string
		var count int
		var title string
		err = rows.Scan(&url, &title, &count)
		if err != nil {
			log.Fatal("searchwild error ", err)
		}

		freq[url] = webInfo{
			title,
			count,
		}
	}

	return freq
}

func (idx dbSql) Search(term string) map[string]webInfo {

	var termID int64
	var freq = map[string]webInfo{}

	err := idx.db.QueryRow(`
	SELECT id FROM terms WHERE name = ?;
	`, term).Scan(&termID)

	if err != nil {
		return freq
	}

	rows, err := idx.db.Query(`
	SELECT urls.name, urls.title, results.term_count FROM urls
	INNER JOIN results ON urls.id = results.url_id AND results.term_id = ?;`, termID)

	if err != nil {
		log.Fatal("search query error ", err)
	}

	for rows.Next() {
		var url string
		var count int
		var title string
		err = rows.Scan(&url, &title, &count)
		if err != nil {
			log.Fatal("search error ", err)
		}

		freq[url] = webInfo{
			title,
			count,
		}
	}

	return freq
}

func lastID(r sql.Result) int64 {
	id, err := r.LastInsertId()
	if err != nil {
		log.Fatal("lasID err ", err)
	}
	return id
}

func (idx dbSql) SearchImg(term string) []imgInfo {

	var freq = []imgInfo{}
	stem, err := snowball.Stem(term, "english", true)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := idx.db.Query(`
	SELECT url, alt, src, title FROM images WHERE alt LIKE '%' || ? || '%';
	`, stem)

	if err != nil {
		log.Fatal("searchImg error ", err)
	}

	for rows.Next() {
		var url string
		var alt string
		var src string
		var title string
		err = rows.Scan(&url, &alt, &src, &title)
		if err != nil {
			log.Fatal("searchimg rows.next error ", err)
		}

		freq = append(freq, imgInfo{url, alt, src, title})
	}

	return freq
}

func insertUrl(u string, title string, wc int, d *dbSql) int64 {
	res, err := d.db.Exec(`
	INSERT INTO urls(name, count, title) VALUES(?,?,?);
	`, u, wc, title)

	if err == nil {
		return lastID(res)
	}

	var urlID int64

	err = d.db.QueryRow(`
		SELECT id FROM urls WHERE name = ?;
		`, u).Scan(urlID)

	if err != nil {
		log.Fatal("insertURL query row: ", err)
	}

	return urlID

}

func insertResult(t int64, u int64, d *dbSql) {
	var resultID int64

	err := d.db.QueryRow(`
		SELECT id FROM results WHERE term_id = ? AND url_id = ?;
		`, t, u).Scan(&resultID)

	if err == nil {
		d.db.QueryRow(`
			UPDATE results SET term_count = term_count + 1 WHERE id = ?;
			`, resultID)
	} else {

		_, err := d.db.Exec(`
			INSERT INTO results(term_id, url_id, term_count) VALUES(?,?,?);
			`, t, u, 1)
		if err != nil {
			log.Fatal("insertResult Queryrow: ", err)
		}
	}

}

func wordStemmer(l []string) []string {

	var stemList []string
	for _, word := range l {
		stem, err := snowball.Stem(word, "english", true)
		if err != nil {
			fmt.Println(err)
		}

		stemList = append(stemList, stem)
	}

	return stemList

}

func (idx dbSql) BiIndexer(wordlist []string, stopwords []string, urlID int64, termIdMap map[string]int64) {

	for i, word := range wordlist {

		if IsStopword(word, stopwords) {
			continue
		}

		if i+1 < len(wordlist) {
			if !IsStopword(wordlist[i+1], stopwords) {

				term1_id := termIdMap[word]
				term2_id := termIdMap[wordlist[i+1]]
				res, err := idx.db.Exec(`
					INSERT INTO bigrams(term1, term2, url_id, freq_count) VALUES(?,?,?,?);
					`, term1_id, term2_id, urlID, 1)

				var biID int64

				if err == nil {
					biID = lastID(res)
				} else {
					err := idx.db.QueryRow(`
						SELECT id FROM bigrams WHERE term1 = ? AND term2 = ?;
						`, term1_id, term2_id).Scan(&biID)

					if err == nil {
						idx.db.QueryRow(`
							UPDATE bigrams SET freq_count = freq_count + 1 WHERE id = ?;
							`, biID)
					}
				}
			}

		}

	}

}

func (idx dbSql) Indexer(extractOut chan extrOut, stopwords []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range extractOut {

		wordlist := wordStemmer(data.Words)
		urlID := insertUrl(data.link, data.Title, len(wordlist), &idx)
		var termIdMap = map[string]int64{}

		for _, word := range wordlist {

			if IsStopword(word, stopwords) {
				continue
			}

			res, err := idx.db.Exec(`
		INSERT INTO terms(name) VALUES(?);
		`, word)

			var termID int64

			if err == nil {
				termID = lastID(res)
			} else {
				err := idx.db.QueryRow(`
			SELECT id FROM terms WHERE name = ?;
			`, word).Scan(&termID)

				if err != nil {
					log.Fatal("Query Indexer: ", err)
				}
			}
			termIdMap[word] = termID
			insertResult(termID, urlID, &idx)
		}

		for _, image := range data.Images {
			_, err := idx.db.Exec(`
		INSERT INTO images(url,alt,src,title) VALUES(?,?,?,?);
		`, data.link, image.alt, image.src, data.Title)
			if err != nil {
				log.Println("adding images err ", err)
			}
		}

		idx.BiIndexer(wordlist, stopwords, urlID, termIdMap)
	}
}

func (idx dbSql) WordsInDoc(url string) int {
	var count int
	row := idx.db.QueryRow(`
	SELECT count FROM urls WHERE name = ?;
	`, url)

	err := row.Scan(&count)
	if err != nil {
		log.Fatal("WordsInDoc err ", err)
	}
	return count

}

func (idx dbSql) DocsInCorp() int {

	mu.Lock()
	defer mu.Unlock()

	var docCount int
	rows, err := idx.db.Query(`
	SELECT COUNT() FROM urls;
	`)

	if err != nil {
		log.Fatal("DocsInCorp err")
	} else {
		for rows.Next() {
			rows.Scan(&docCount)
		}
	}

	return docCount

}
