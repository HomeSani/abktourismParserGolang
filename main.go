package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetDoc(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln(err)
	}

	return doc
}

func GetMaxPageIndex() int {
	maxPageIndex := 1

	for true {
		url := fmt.Sprintf("http://b2b.abktourism.kz/search_tour?TOWNFROMINC=10&STATEINC=6&CHECKIN_BEG=20230418&NIGHTS_FROM=7&CHECKIN_END=20230419&NIGHTS_TILL=9&ADULT=2&CURRENCY=4&CHILD=0&TOWNS_ANY=1&STARS_ANY=1&HOTELS_ANY=1&MEALS_ANY=1&PRICEPAGE=%d&DOLOAD=1", maxPageIndex)
		doc := GetDoc(url)

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
	// maxPageIndex := getMaxPageIndex()

	for i := 1; i <= 1; i++ {
		url := fmt.Sprintf("http://b2b.abktourism.kz/search_tour?TOWNFROMINC=10&STATEINC=6&CHECKIN_BEG=20230418&NIGHTS_FROM=7&CHECKIN_END=20230419&NIGHTS_TILL=9&ADULT=2&CURRENCY=4&CHILD=0&TOWNS_ANY=1&STARS_ANY=1&HOTELS_ANY=1&MEALS_ANY=1&PRICEPAGE=%d&DOLOAD=1", i)
		doc := GetDoc(url)

		doc.Find("tr.stats").Each(func(i int, s *goquery.Selection) {
			tourStartTime := strings.TrimSpace(s.Find("td.sortie").Text())
			tourName := strings.TrimSpace(s.Find("td.tour").Text())
			nigthsCount := strings.TrimSpace(s.Find("td.c").First().Text())
			hotelName := strings.TrimSpace(s.Find("td.link-hotel").Text())
			havePlacesTmp, _ := s.Find("td.nw").Children().Attr("title")
			havePlaces := strings.TrimSpace(havePlacesTmp)
			if havePlaces == "" {
				havePlaces = "да"
			}
			nutrition := strings.TrimSpace(s.Find("td.nw").Next().Text())
			roomAndAccommodation := strings.TrimSpace(s.Find("td.nw").Next().Next().Text())
			price := strings.TrimSpace(s.Find("td.price").Children().Text())
			priceType := strings.TrimSpace(s.Find("td.type_price").Children().Text())

			fmt.Println("===========================================================================================")
			fmt.Println(tourStartTime, tourName, nigthsCount, hotelName, havePlaces, nutrition, roomAndAccommodation, price, priceType)
			fmt.Println("===========================================================================================")
		})
	}
}
