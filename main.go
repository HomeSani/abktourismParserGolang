package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
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

	db, err := sql.Open("mysql", "lol:lool#@tcp(localhost:3306)/tours_db")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	for i := 1; i <= 1; i++ {
		url := fmt.Sprintf("http://b2b.abktourism.kz/search_tour?TOWNFROMINC=10&STATEINC=6&CHECKIN_BEG=20230418&NIGHTS_FROM=7&CHECKIN_END=20230419&NIGHTS_TILL=9&ADULT=2&CURRENCY=4&CHILD=0&TOWNS_ANY=1&STARS_ANY=1&HOTELS_ANY=1&MEALS_ANY=1&PRICEPAGE=%d&DOLOAD=1", i)
		doc := GetDoc(url)

		doc.Find("tr.stats").Each(func(i int, s *goquery.Selection) {
			s.Find("td.sortie").RemoveClass("transp_icon_1")
			s.Find("td.price span").RemoveClass("bron price_button")
			tourStartDate := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s.Find("td.sortie").Text(), "\n", " "), " ", ""))
			tourName := strings.TrimSpace(s.Find("td.tour").Text())
			nigthsCountString := strings.TrimSpace(s.Find("td.c").First().Text())
			nigthsCount, _ := strconv.Atoi(nigthsCountString)
			hotelName := strings.TrimSpace(s.Find("td.link-hotel").Text())
			havePlacesTmp, _ := s.Find("td.nw").Children().Attr("title")
			havePlaces := strings.TrimSpace(havePlacesTmp)
			if havePlaces == "" {
				havePlaces = "да"
			}
			nutrition := strings.TrimSpace(s.Find("td.link-hotel").Next().Next().Text())
			roomAndAccommodation := strings.TrimSpace(s.Find("td.link-hotel").Next().Next().Next().Text())
			price := strings.ReplaceAll(s.Find("td.td_price .price").Text(), " ", "")
			price = strings.ReplaceAll(price, "\n", "")
			priceType := strings.TrimSpace(s.Find("td.type_price").Children().Text())

			fmt.Println("===========================================================================================")
			fmt.Println(tourStartDate, tourName, nigthsCount, hotelName, havePlaces, nutrition, roomAndAccommodation, price, priceType)
			fmt.Println("===========================================================================================")

			stmt, err := db.Prepare("INSERT INTO tours(start_time, nights_count, tour_name, hotel, avalible_places, nutrition, room_type, price, price_type) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Panicln(err)
			}

			defer stmt.Close()

			result, err := stmt.Exec(tourStartDate, nigthsCount, tourName, hotelName, havePlaces, nutrition, roomAndAccommodation, price, priceType)
			if err != nil {
				log.Panicln(err)
			}
			fmt.Println(result.LastInsertId())
		})
	}
}
