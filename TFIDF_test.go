package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kljensen/snowball"
)

var in = NewDb()

func TestCase1(t *testing.T) {

	var test1 = []result{

		{Term: "comput scienc", Url: "https://www.ucsc.edu/azindex/", TFIDF: 0.014249790444258174},
	}

	term, err := snowball.Stem("computer", "english", true)
	if err != nil {
		fmt.Println(err)
	}
	term2, err := snowball.Stem("science", "english", true)
	if err != nil {
		fmt.Println(err)
	}
	isBigram := true
	isWild := true
	go DbSearch(&in)

	got := DbCalculateTFIDF(term, &in, isBigram, term2, isWild)

	testcase := reflect.DeepEqual(got, test1)

	if !testcase {
		for i := range got {
			// t.Errorf("{Term: \"%s\", Url: \"%s\", TFIDF: %v},\n", got[i].Term, got[i].Url, got[i].TFIDF)
			if test1[i].TFIDF != got[i].TFIDF {

				t.Errorf("\nExpected: %v\ngot: %v\n", test1[i], got[i])
			}
		}
	}
}

func TestCase2(t *testing.T) {

	var test2 = []result{

		{Term: "class", Url: "https://www.ucsc.edu/academics/", TFIDF: 0.004757373929590866},

		{Term: "class", Url: "https://www.ucsc.edu/research/", TFIDF: 0.002849002849002849},

		{Term: "class", Url: "https://www.ucsc.edu/admissions/", TFIDF: 0.002688172043010753},

		{Term: "class", Url: "https://www.ucsc.edu/author/lmnielseucsc-edu/", TFIDF: 0.0025575447570332483},

		{Term: "class", Url: "https://www.ucsc.edu/author/milpowelucsc-edu/", TFIDF: 0.0025575447570332483},

		{Term: "class", Url: "https://www.ucsc.edu/author/gwenjucsc-edu/", TFIDF: 0.0025575447570332483},

		{Term: "class", Url: "https://www.ucsc.edu/author/raknightucsc-edu/", TFIDF: 0.002551020408163265},

		{Term: "class", Url: "https://www.ucsc.edu/research/undergraduate-research/", TFIDF: 0.002551020408163265},

		{Term: "class", Url: "https://www.ucsc.edu/campus/campus-galleries-and-theaters/", TFIDF: 0.002531645569620253},

		{Term: "class", Url: "https://www.ucsc.edu/search/", TFIDF: 0.0025252525252525255},

		{Term: "class", Url: "https://www.ucsc.edu/programs-and-units/", TFIDF: 0.002296211251435132},

		{Term: "class", Url: "https://www.ucsc.edu/about/", TFIDF: 0.002285714285714286},

		{Term: "class", Url: "https://www.ucsc.edu/programs-and-units/", TFIDF: 0.0022222222222222222},

		{Term: "class", Url: "https://www.ucsc.edu/campus/", TFIDF: 0.002150537634408602},

		{Term: "class", Url: "https://www.ucsc.edu/", TFIDF: 0.0021119324181626186},

		{Term: "class", Url: "https://www.ucsc.edu", TFIDF: 0.0021119324181626186},

		{Term: "class", Url: "https://www.ucsc.edu/azindex/", TFIDF: 0.0020955574182732607},

		{Term: "class", Url: "https://www.ucsc.edu/feedback/", TFIDF: 0.0020408163265306124},

		{Term: "class", Url: "https://www.ucsc.edu/people/", TFIDF: 0.001996007984031936},

		{Term: "class", Url: "https://www.ucsc.edu/2023-highlights/", TFIDF: 0.0018315018315018315},

		{Term: "class", Url: "https://www.ucsc.edu/residential-colleges/", TFIDF: 0.0017953321364452424},

		{Term: "class", Url: "https://www.ucsc.edu/address-and-phone/", TFIDF: 0.0016420361247947454},

		{Term: "class", Url: "https://www.ucsc.edu/campus/visit/maps-directions/", TFIDF: 0.0016207455429497568},

		{Term: "class", Url: "https://www.ucsc.edu/land-acknowledgment/", TFIDF: 0.001455604075691412},

		{Term: "class", Url: "https://www.ucsc.edu/campus/mascot/", TFIDF: 0.0014104372355430183},

		{Term: "class", Url: "https://www.ucsc.edu/mission-and-vision/", TFIDF: 0.0013908205841446453},

		{Term: "class", Url: "https://www.ucsc.edu/about/achievements-facts-and-figures/", TFIDF: 0.0013679890560875513},

		{Term: "class", Url: "https://www.ucsc.edu/principles-community/", TFIDF: 0.0013280212483399733},

		{Term: "class", Url: "https://www.ucsc.edu/about/leadership/", TFIDF: 0.0013175230566534915},

		{Term: "class", Url: "https://www.ucsc.edu/better-together/", TFIDF: 0.0013157894736842105},

		{Term: "class", Url: "https://www.ucsc.edu/campus/visit/", TFIDF: 0.00129366106080207},

		{Term: "class", Url: "https://www.ucsc.edu/about/overview/", TFIDF: 0.0012345679012345679},

		{Term: "class", Url: "https://www.ucsc.edu/campus-destinations/", TFIDF: 0.0010626992561105207},

		{Term: "class", Url: "https://www.ucsc.edu/privacy-policy/", TFIDF: 0.0009861932938856016},
	}
	term, err := snowball.Stem("class", "english", true)
	if err != nil {
		fmt.Println(err)
	}
	term2 := ""
	isBigram := false
	isWild := true

	got := DbCalculateTFIDF(term, &in, isBigram, term2, isWild)

	secondCase := reflect.DeepEqual(got, test2)

	if !secondCase {
		for i := range got {
			// t.Errorf("\n{Term: \"%s\", Url: \"%s\", TFIDF: %v},\n", got[i].Term, got[i].Url, got[i].TFIDF)
			if test2[i].TFIDF != got[i].TFIDF {
				t.Errorf("\nExpected: %v\ngot: %v\n", test2[i], got[i])
			}
		}
	}
}
