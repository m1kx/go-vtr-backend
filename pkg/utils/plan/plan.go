package plan

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1kx/go-vtr-backend/pkg/utils"
)

func format_row(row []string) {
	split := strings.Split(row[2], " ")
	if len(split) > 2 {
		row[2] = fmt.Sprintf("%s %s", split[0], split[2])
	}
}

func Scrape(day string) (data [][]string, base string, wd string, err error) {
	res, err := http.Get(fmt.Sprintf("https://lmg-anrath.de/aktuelle_plaene/Vertretungsplan/%s/subst_001.htm", day))
	if err != nil {
		return nil, "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, "", "", err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", "", err
	}

	// date handling
	heading := doc.Find(".mon_title").Text()
	heading_fmt := strings.ReplaceAll(strings.Split(heading, " ")[0], ".", "-")
	date, err := time.Parse("2-1-2006", heading_fmt)
	if err != nil {
		return nil, "", "", err
	}
	date_string := date.Format("02-01-2006")
	_ = date_string
	curr_time := time.Now()
	if day == "morgen" {
		curr_time = curr_time.AddDate(0, 0, 1)
	}
	curr_date := curr_time.Format("02-01-2006")
	_ = curr_date

	weekday := ""
	switch date.Weekday() {
	case time.Sunday:
		weekday = "Sonntag"
	case time.Monday:
		weekday = "Montag"
	case time.Tuesday:
		weekday = "Dienstag"
	case time.Wednesday:
		weekday = "Mittwoch"
	case time.Thursday:
		weekday = "Donnerstag"
	case time.Friday:
		weekday = "Freitag"
	case time.Saturday:
		weekday = "Samstag"
	}

	if curr_time.Weekday() == date.Weekday() {
		weekday = fmt.Sprintf("%s (%s)", weekday, day)
	}

	var rows []*goquery.Selection
	doc.Find(".list").Each(func(i int, s *goquery.Selection) {
		rows = append(rows, s)
	})

	var rows_formatted [][]string
	for i := 9; i < len(rows); i++ {
		var this_row []string
		for x := 0; x < 7; x++ {
			this_row = append(this_row, rows[i+x].Text())
		}
		i += 7
		rows_formatted = append(rows_formatted, this_row)
	}

	for i := 0; i < len(rows_formatted); i++ {
		format_row(rows_formatted[i])
	}

	data_string := ""
	for i := 0; i < len(rows_formatted); i++ {
		for x := 0; x < len(rows_formatted[i]); x++ {
			data_string = fmt.Sprintf("%s%s", data_string, rows_formatted[i][x])
		}
		data_string = fmt.Sprintf("%s%s", data_string, ":")
	}

	base64_string := utils.EncodeBase64(data_string)

	return rows_formatted, base64_string, weekday, nil
}
