package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func getMaxPageIndex() int {
	maxPageIndex := 1

	for true {
		url := fmt.Sprintf("http://b2b.abktourism.kz/search_tour?TOWNFROMINC=10&STATEINC=6&CHECKIN_BEG=20230418&NIGHTS_FROM=7&CHECKIN_END=20230419&NIGHTS_TILL=9&ADULT=2&CURRENCY=4&CHILD=0&TOWNS_ANY=1&STARS_ANY=1&HOTELS_ANY=1&MEALS_ANY=1&PRICEPAGE=%d&DOLOAD=1", maxPageIndex)

		res, err := http.Get(url)
		if err != nil {
			log.Panicln(err)
		}

		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Panicln(err)
		}

		lastPageIndex, err := strconv.Atoi(doc.Find("span.page").Last().Text())
		if err != nil {
			log.Panicln(err)
		}

		if maxPageIndex < lastPageIndex {
			maxPageIndex = lastPageIndex
		} else if maxPageIndex == lastPageIndex+1 {
			break
		}
	}

	return maxPageIndex
}

func main() {
	maxPageIndex := getMaxPageIndex()

	fmt.Println(maxPageIndex)
}
